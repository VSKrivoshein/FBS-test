package grpc_api_test

import (
	context "context"
	api "github.com/VSKrivoshein/FBS-test/internal/app/api/grpc_api/proto"
	"github.com/VSKrivoshein/FBS-test/internal/app/configs"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"os"
	"testing"
)

const bufSize = 1024 * 1024

var listen *bufconn.Listener
var TestGrpcClient api.FibonacciClient

func TestMain(m *testing.M) {
	listen = bufconn.Listen(bufSize)
	service := configs.InitService()
	go func() {
		if err := configs.StartGrpc(service, listen); err != nil {
			log.Fatal(err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure(),
	)


	if err != nil {
		os.Exit(1)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Fatal closing resp conn: %v", err.Error())
		}
	}()

	TestGrpcClient = api.NewFibonacciClient(conn)
	exitVal := m.Run()
	os.Exit(exitVal)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return listen.Dial()
}