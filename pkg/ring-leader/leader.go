package ringLeader

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	omap "github.com/elliotchance/orderedmap/v2"
	"github.com/kolharsam/go-delta/pkg/config"
	pb "github.com/kolharsam/go-delta/pkg/grpc"
	"github.com/kolharsam/go-delta/pkg/lib"
)

type workerId = string

type taskWorkerInfo struct {
	ServiceId     string    `json:"service_id"`
	ServiceHost   string    `json:"service_host"`
	Port          uint32    `json:"port"`
	LastHeartBeat time.Time `json:"last_heartbeat"`
}

type taskWorkers struct {
	mtx     sync.RWMutex
	workers *omap.OrderedMap[workerId, *taskWorkerInfo]
	next    uint32
}

// nextWorkerForTask distributes the tasks amongst the connected workers
// by using the round-robin algorithm. The omap isn't the best data structure
// to get this done efficiently. Improvements will be made in the future versions
func nextWorkerForTask(ts *taskWorkers) *taskWorkerInfo {
	ts.mtx.Lock()
	n := atomic.AddUint32(&ts.next, 1)
	ts.mtx.Unlock()

	if int(n) > ts.workers.Len() {
		ts.mtx.Lock()
		atomic.StoreUint32(&ts.next, 1)
		ts.mtx.Unlock()
		n = 1
	}

	nxtWorkerIndex := (int(n) - 1) % ts.workers.Len()

	iter := 0
	for el := ts.workers.Front(); el != nil; el = el.Next() {
		if iter == nxtWorkerIndex {
			return el.Value
		}
		iter++
	}
	return nil
}

type ringLeaderServer struct {
	pb.UnimplementedRingLeaderServer
	activeServers *taskWorkers
	logger        *zap.Logger
	leaderHost    string
	leaderPort    uint32
	appConfig     *config.DeltaConfig
}

type connectionRequest struct {
	serviceId   string
	serviceHost string
	port        uint32
	timeStamp   string
}

func (tsi *taskWorkerInfo) updateHeartbeatTimestamp(timeStamp string) error {
	tm, err := time.Parse(time.RFC3339, timeStamp)
	if err != nil {
		return err
	}
	tsi.LastHeartBeat = tm
	return nil
}

func (ts *taskWorkers) addNewService(connectRequest connectionRequest) error {
	tm, err := time.Parse(time.RFC3339, connectRequest.timeStamp)

	if err != nil {
		return err
	}

	tsi := taskWorkerInfo{
		ServiceId:     connectRequest.serviceId,
		ServiceHost:   connectRequest.serviceHost,
		Port:          connectRequest.port,
		LastHeartBeat: tm,
	}

	ts.workers.Set(connectRequest.serviceId, &tsi)

	return nil
}

func newWorkerServiceClient(host string, port uint32) (pb.WorkerClient, error) {
	workerTarget := fmt.Sprintf("%s:%d", host, port)

	conn, err := grpc.NewClient(workerTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return pb.NewWorkerClient(conn), nil
}

func (ts *taskWorkers) updateServiceHeartbeat(serviceId, timestamp string) error {
	if taskWorker, ok := ts.workers.Get(serviceId); ok {
		if err := taskWorker.updateHeartbeatTimestamp(timestamp); err != nil {
			return err
		}
	}
	return nil
}

func (ts *taskWorkers) removeService(serviceId string) *taskWorkerInfo {
	val, ok := ts.workers.Get(serviceId)
	if !ok {
		return nil
	}
	ts.workers.Delete(serviceId)
	return val
}

func (rls *ringLeaderServer) Hearbeat(stream grpc.BidiStreamingServer[pb.HeartbeatFromWorker, pb.HeartbeatFromLeader]) error {
	for {
		beat, err := stream.Recv()
		if err == io.EOF {
			rls.activeServers.mtx.Lock()
			rls.activeServers.removeService(beat.ServiceId)
			rls.activeServers.mtx.Unlock()
			return nil
		}

		if err != nil {
			rls.logger.Error("failed to receive heartbeat from worker...")
			return err
		}

		workerId := beat.ServiceId
		beatTime := beat.Timestamp.AsTime().Format(time.RFC3339)

		rls.activeServers.mtx.Lock()
		rls.activeServers.updateServiceHeartbeat(workerId, beatTime)
		rls.activeServers.mtx.Unlock()

		rls.logger.Info("updated the worker status from heartbeat...",
			zap.String("worker_id", workerId),
		)

		stream.Send(&pb.HeartbeatFromLeader{
			Timestamp: timestamppb.Now(),
		})
	}
}

func (rls *ringLeaderServer) CheckHearbeats() {
	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for range ticker.C {
		if rls.activeServers.workers.Len() == 0 {
			continue
		}

		rls.activeServers.mtx.RLock()

		for el := rls.activeServers.workers.Front(); el != nil; el = el.Next() {
			if time.Since(el.Value.LastHeartBeat) >= (time.Second * 15) {
				rls.logger.Warn("worker seems to be down...",
					zap.String("worker_id", el.Value.ServiceId))
			}
		}

		rls.activeServers.mtx.RUnlock()

	}
}

func (rls *ringLeaderServer) Connect(ctx context.Context, connReq *pb.ConnectRequest) (*pb.ConnectAck, error) {
	connectRequest := connectionRequest{
		serviceId:   connReq.GetServiceId(),
		serviceHost: connReq.GetServiceHost(),
		port:        connReq.GetPort(),
		timeStamp:   connReq.GetTimestamp().String(),
	}

	rls.activeServers.mtx.Lock()
	err := rls.activeServers.addNewService(connectRequest)
	rls.activeServers.mtx.Unlock()

	if err != nil {
		return nil, err
	}

	rls.logger.Info("connected with new worker...",
		zap.String("worker_host", connectRequest.serviceHost),
		zap.Uint32("worker_port", connectRequest.port),
		zap.String("worker_id", connectRequest.serviceId),
	)

	return &pb.ConnectAck{
		Host:      rls.leaderHost,
		Port:      rls.leaderPort,
		Timestamp: timestamppb.Now(),
	}, nil
}

func newServer(host string, port uint32, logger *zap.Logger, config *config.DeltaConfig) *ringLeaderServer {
	s := &ringLeaderServer{
		activeServers: &taskWorkers{
			workers: omap.NewOrderedMap[string, *taskWorkerInfo](),
		},
		logger:     logger,
		leaderHost: host,
		leaderPort: port,
		appConfig:  config,
	}
	return s
}

func GetListenerAndServer(host string, port uint32, config *config.DeltaConfig) (net.Listener, *grpc.Server, *ringLeaderServer, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, nil, nil, err
	}

	logger, err := lib.GetLogger()

	if err != nil {
		log.Fatalf("failed to initiate logger for ring-leader [%v]", err)
		return nil, nil, nil, err
	}

	grpcServer := grpc.NewServer()
	serverCtx := newServer(host, port, logger, config)
	pb.RegisterRingLeaderServer(grpcServer, serverCtx)
	return listener, grpcServer, serverCtx, nil
}
