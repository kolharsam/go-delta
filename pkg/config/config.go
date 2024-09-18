package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type BloomFilterConfig struct {
	Bits             uint64 `json:"bits" toml:"bits"`
	NumHashFunctions uint   `json:"num_hash_functions" toml:"num_hash_functions"`
}

type RingLeaderConfig struct {
	Connections ConnectionsConfig `json:"connections" toml:"connections"`
}

type ConnectionsConfig struct {
	MaxRetries         int `json:"max_retries" toml:"max_retries"`
	TimeBetweenRetries int `json:"time_between_retries" toml:"time_between_retries"`
}

type WorkerConfig struct {
	HeartbeatInterval int               `json:"heartbeat_interval" toml:"heartbeat_interval"`
	BackoffMax        int               `json:"backoff_max" toml:"backoff_max"`
	Connections       ConnectionsConfig `json:"connections" toml:"connections"`
}

type DeltaConfig struct {
	Title             string            `json:"title" toml:"title"`
	RingLeaderConfig  RingLeaderConfig  `json:"ring_leader" toml:"ring-leader"`
	WorkerConfig      WorkerConfig      `json:"worker" toml:"worker"`
	BloomFilterConfig BloomFilterConfig `json:"bloom" toml:"bloom"`
}

var (
	defaultConfig = &DeltaConfig{
		Title: "go-delta",
		RingLeaderConfig: RingLeaderConfig{
			Connections: ConnectionsConfig{
				MaxRetries:         10,
				TimeBetweenRetries: 5, // NOTE: this is in seconds
			},
		},
		WorkerConfig: WorkerConfig{
			HeartbeatInterval: 2,
			BackoffMax:        2,
			Connections: ConnectionsConfig{
				MaxRetries:         10,
				TimeBetweenRetries: 5,
			},
		},
		BloomFilterConfig: BloomFilterConfig{
			Bits:             1000,
			NumHashFunctions: 3,
		},
	}
)

func ParseConfig(fileName string) (*DeltaConfig, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return defaultConfig, err
	}

	var config DeltaConfig

	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
