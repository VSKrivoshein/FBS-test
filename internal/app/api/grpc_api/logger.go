package grpc_api

import (
	context "context"
	"github.com/VSKrivoshein/FBS-test/internal/app/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

func Logger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	// Start timer
	start := time.Now()

	res, err := handler(ctx, req)

	// Stop timer
	duration := utils.GetDurationInMilliseconds(start)

	entry := log.WithFields(log.Fields{
		"method":   info.FullMethod,
		"duration": duration,
	})

	if err != nil {
		entry.Error(err)
		return res, err
	}
	entry.Info("success")
	return res, err
}
