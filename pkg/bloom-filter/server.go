package bloomfilter

import (
	"fmt"
	"log"
	"net"

	"github.com/kolharsam/go-delta/pkg/config"
	"github.com/kolharsam/go-delta/pkg/lib"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/kolharsam/go-delta/pkg/grpc"
)

type bloomFilterServerCtx struct {
	pb.BloomFilterServer
	logger    *zap.Logger
	appConfig *config.DeltaConfig
}

func newServerCtx(logger *zap.Logger, config *config.DeltaConfig) *bloomFilterServerCtx {
	s := &bloomFilterServerCtx{
		logger:    logger,
		appConfig: config,
	}
	return s
}

func GetListenerAndServer(host string, port uint32, config *config.DeltaConfig) (net.Listener, *grpc.Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, nil, err
	}

	logger, err := lib.GetLogger()

	if err != nil {
		log.Fatalf("failed to initiate logger for bloom-filter service [%v]", err)
		return nil, nil, err
	}

	grpcServer := grpc.NewServer()
	serverCtx := newServerCtx(logger, config)
	pb.RegisterBloomFilterServer(grpcServer, serverCtx)

	return listener, grpcServer, nil
}
