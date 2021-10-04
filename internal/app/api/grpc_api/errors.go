package grpc_api

import (
	context "context"
	e "github.com/VSKrivoshein/FBS-test/internal/app/custom_err"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Err(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	res, err := handler(ctx, req)

	if err == nil {
		return res, err
	}

	unwrappedErr := e.UnwrapRecursive(err)
	customErr, ok := unwrappedErr.(*e.CustomErrorWithCode)
	if !ok {
		logrus.Fatalf("error: error was not found in unwrappedErr.(*e.CustomErrorWithCode)")
	}

	return res, status.Error(customErr.GrpcCode, customErr.ErrorForUser())
}
