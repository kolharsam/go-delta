package main

import (
	"flag"
	"log"

	bloomfilter "github.com/kolharsam/go-delta/pkg/bloom-filter"
	"github.com/kolharsam/go-delta/pkg/config"
)

var (
	host       = flag.String("host", "localhost", "--host \"localhost\"")
	port       = flag.Uint("port", 8082, "--port 8082")
	configFile = flag.String("config", "config.toml", "--config ../../<path-to-config>")
)

func main() {
	flag.Parse()

	appConfig, err := config.ParseConfig(*configFile)

	if err != nil && appConfig != nil {
		log.Println("config file not provided...switching to defaults...")
	} else if err != nil && appConfig == nil {
		log.Fatalf("there's an issue with the config file provided...[%v]", err)
	}

	log.Println("applied config successfully...")

	lis, server, err := bloomfilter.GetListenerAndServer(*host, uint32(*port), appConfig)
	if err != nil {
		log.Fatalf("failed to setup bloom-filter server %v", err)
	}

	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("failure at bloom-filter server at [%s:%d]", *host, *port)
	}
}
