package configs

import (
	"fmt"
	"github.com/VSKrivoshein/FBS-test/internal/app/api/grpc_api"
	api "github.com/VSKrivoshein/FBS-test/internal/app/api/grpc_api/proto"
	"github.com/VSKrivoshein/FBS-test/internal/app/api/rest_api"
	"github.com/VSKrivoshein/FBS-test/internal/app/cache"
	"github.com/VSKrivoshein/FBS-test/internal/app/services/fiboncci"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func InitService() fiboncci.Service {
	rdb, err := cache.NewRdb()
	if err != nil {
		log.Fatal(err)
	}
	if err := rdb.PrepareRDB(); err != nil {
		log.Fatal(err)
	}
	return fiboncci.NewService(rdb)
}

func StartRest(service fiboncci.Service) error {
	restPort := fmt.Sprintf(":%s", os.Getenv("REST_PORT"))
	log.Infof("Starting rest server on port %s", restPort)
	handler := rest_api.NewHandler(service)
	return http.ListenAndServe(restPort, handler.InitRoutes())
}

func StartGrpc(service fiboncci.Service, listen net.Listener) error {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpc_api.Err, grpc_api.Logger),
	)

	srv := grpc_api.NewGRPCServer(service)
	api.RegisterFibonacciServer(s, srv)

	return s.Serve(listen)
}

func GracefulShutdown() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	return fmt.Errorf("%s", <-signals)
}