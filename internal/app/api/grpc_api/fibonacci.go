package grpc_api

import (
	context "context"
	"fmt"
	api "github.com/VSKrivoshein/FBS-test/internal/app/api/grpc_api/proto"
	e "github.com/VSKrivoshein/FBS-test/internal/app/custom_err"
	"github.com/VSKrivoshein/FBS-test/internal/app/services/fiboncci"
)

type GRPCServer struct {
	Service fiboncci.Service
	api.UnimplementedFibonacciServer
}

func NewGRPCServer(service fiboncci.Service) *GRPCServer {
	return &GRPCServer{Service: service}
}

func (s *GRPCServer) CalcFibonacciSequence(ctx context.Context, req *api.CalcFibonacciSequenceReq) (*api.CalcFibonacciSequenceResponse, error) {
	fibonacciSequence, err := s.Service.Calculate(req.GetX(), req.GetY())
	if err != nil {
		return nil, fmt.Errorf(e.GetInfo(), err)
	}

	return &api.CalcFibonacciSequenceResponse{FibonacciSequence: fibonacciSequence}, nil
}
