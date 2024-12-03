package main

import (
	"flag"

	cmd "github.com/kolharsam/go-delta/pkg/cli"
)

var (
	port = flag.Uint("port", 8082, "--port 8082")
)

func main() {
	cmd.Execute(*port)
}
