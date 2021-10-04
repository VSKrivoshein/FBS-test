package grpc_api_test

import (
	context "context"
	api "github.com/VSKrivoshein/FBS-test/internal/app/api/grpc_api/proto"
	e "github.com/VSKrivoshein/FBS-test/internal/app/custom_err"
	"github.com/VSKrivoshein/FBS-test/internal/app/services/fiboncci"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestFibonacciGrpc(t *testing.T) {
	tests := []struct {
		name            string
		x               uint32
		y               uint32
		isFail          bool
		wontCode        codes.Code
		wontBodySuccess []string
		wontErrorMsg    string
	}{
		{
			name:            "zero zero",
			x:               0,
			y:               0,
			isFail:          false,
			wontCode:        codes.OK,
			wontBodySuccess: []string{"0"},
			wontErrorMsg:    "",
		},
		{
			name:            "part of sequences in cash",
			x:               0,
			y:               1,
			isFail:          false,
			wontCode:        codes.OK,
			wontBodySuccess: []string{"0", "1"},
			wontErrorMsg:    "",
		},
		{
			name:            "calculate new values with part cache",
			x:               0,
			y:               5,
			isFail:          false,
			wontCode:        codes.OK,
			wontBodySuccess: []string{"0", "1", "1", "2", "3", "5"},
			wontErrorMsg:    "",
		},
		{
			name:            "calculate new values without cache",
			x:               7,
			y:               9,
			isFail:          false,
			wontCode:        codes.OK,
			wontBodySuccess: []string{"13", "21", "34"},
			wontErrorMsg:    "",
		},
		{
			name:            "second value less than first",
			x:               2,
			y:               1,
			isFail:          true,
			wontCode:        codes.InvalidArgument,
			wontBodySuccess: nil,
			wontErrorMsg:    e.ErrSecondNumberLessThanFirst.Error(),
		},
		{
			name:            "max value of exceeded",
			x:               0,
			y:               fiboncci.MaxIntFibonacci + 1,
			isFail:          true,
			wontCode:        codes.InvalidArgument,
			wontBodySuccess: nil,
			wontErrorMsg:    e.ErrMaxIntFibonacci.Error(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			resp, err := TestGrpcClient.CalcFibonacciSequence(ctx, &api.CalcFibonacciSequenceReq{
				X: test.x,
				Y: test.y,
			})

			if test.isFail {
				st, ok := status.FromError(err)
				if !ok {
					t.Fatal("is not able to get status from err: status, ok := status.FromError(err)")
				}
				assert.Equal(t, test.wontErrorMsg, st.Message())
				assert.Equal(t, test.wontCode, st.Code())
			}
			if test.isFail {
				st, ok := status.FromError(err)
				if !ok {
					t.Fatal("is not able to get status from err: status, ok := status.FromError(err)")
				}
				assert.Equal(t, test.wontCode, st.Code())
				assert.ElementsMatch(t, resp.GetFibonacciSequence(), test.wontBodySuccess)
			}
		})
	}
}
