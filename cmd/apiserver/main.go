package main

import (
	"fmt"
	"github.com/VSKrivoshein/FBS-test/internal/app/configs"
	"net"
	"os"

	"github.com/VSKrivoshein/FBS-test/internal/app/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	utils.UseJSONLogFormat()
	errors := make(chan error)
	service := configs.InitService()

	go func() {
		errors <- configs.StartRest(service)
	}()

	go func() {
		grpcPort := fmt.Sprintf(":%s", os.Getenv("GRPC_PORT"))
		log.Infof("Starting grpc server on port %s", grpcPort)
		listen, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("Lister start fatal: %v", err)
		}
		errors <- configs.StartGrpc(service, listen)
	}()

	go func() {
		errors <- configs.GracefulShutdown()
	}()

	log.Fatalf("Terminated %s", <-errors)
}
