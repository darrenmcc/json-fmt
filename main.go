package main

import (
	"log"
	"os"

	"github.com/darrenmcc/gizmo"
	"github.com/darrenmcc/json-fmt/api"
)

func main() {
	err := gizmo.Run(api.NewService(mustEnv("API_SECRET")))
	if err != nil {
		log.Fatalf("unable to start service: %s", err)
	}
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic(k + " environment variable not set")
	}
	return v
}
