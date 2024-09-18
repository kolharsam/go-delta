package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
)

var (
	port = flag.String("port", "8080", "--port 8080")
)

func main() {
	err := godotenv.Load()
	flag.Parse()
	if err != nil {
		log.Fatalf("failed to set up environment variables [%v]", err)
	}

	// serverPort := lib.MakePortString(*port)

	// schedulerapi.Run(serverPort)
}
