package rest_api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/VSKrivoshein/FBS-test/internal/app/api/rest_api"
	e "github.com/VSKrivoshein/FBS-test/internal/app/custom_err"
	"github.com/VSKrivoshein/FBS-test/internal/app/services/fiboncci"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	contentType = "application/json"
)

func TestFibonacciRest(t *testing.T) {
	tests := []struct {
		name            string
		x               uint32
		y               uint32
		isFail          bool
		wontCode        int
		wontBodySuccess []string
		wontErrorMsg    string
	}{
		{
			name:            "zero zero",
			x:               0,
			y:               0,
			isFail:          false,
			wontCode:        200,
			wontBodySuccess: []string{"0"},
			wontErrorMsg:    "",
		},
		{
			name:            "part of sequences in cash",
			x:               0,
			y:               1,
			isFail:          false,
			wontCode:        200,
			wontBodySuccess: []string{"0", "1"},
			wontErrorMsg:    "",
		},
		{
			name:            "calculate new values with part cache",
			x:               0,
			y:               5,
			isFail:          false,
			wontCode:        200,
			wontBodySuccess: []string{"0", "1", "1", "2", "3", "5"},
			wontErrorMsg:    "",
		},
		{
			name:            "calculate new values without cache",
			x:               7,
			y:               9,
			isFail:          false,
			wontCode:        200,
			wontBodySuccess: []string{"13", "21", "34"},
			wontErrorMsg:    "",
		},
		{
			name:            "second value less than first",
			x:               2,
			y:               1,
			isFail:          true,
			wontCode:        http.StatusUnprocessableEntity,
			wontBodySuccess: nil,
			wontErrorMsg:    e.ErrSecondNumberLessThanFirst.Error(),
		},
		{
			name:            "max value of exceeded",
			x:               0,
			y:               fiboncci.MaxIntFibonacci + 1,
			isFail:          true,
			wontCode:        http.StatusUnprocessableEntity,
			wontBodySuccess: nil,
			wontErrorMsg:    e.ErrMaxIntFibonacci.Error(),
		},
	}

	url := fmt.Sprintf("%v/fibonacci", TestSrv.URL)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			content, err := json.Marshal(rest_api.FibonacciInput{
				X: test.x,
				Y: test.y,
			})
			if err != nil {
				t.Fatalf("TestFibonacci marshal fatal: %v", err)
			}

			body := bytes.NewBuffer(content)
			resp, err := http.Post(url, contentType, body)
			defer RespClose(t, resp)
			if err != nil {
				t.Fatalf("TestFibonacci marshal fatal: %v", err)
			}

			if test.isFail {
				var respBody struct {
					Error string `json:"error"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
					t.Fatalf("TestFibonacci Fatal json.NewDecoder(resp.Body).Decode(respStruct) %v", err)
				}
				assert.Equal(t, resp.StatusCode, test.wontCode)
				assert.Equal(t, respBody.Error, test.wontErrorMsg)
			}

			if !test.isFail {
				var respBody []string
				if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
					t.Fatalf("TestFibonacci Fatal json.NewDecoder(resp.Body).Decode(respStruct) %v", err)
				}
				assert.Equal(t, resp.StatusCode, test.wontCode)
				assert.ElementsMatch(t, respBody, test.wontBodySuccess)
			}

		})
	}
}

func RespClose(t *testing.T, resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		t.Fatalf("Fatal closing resp")
	}
}
