package rest_api

import (
	e "github.com/VSKrivoshein/FBS-test/internal/app/custom_err"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"net/http"
)

type FibonacciInput struct {
	X uint32 `json:"x,string"`
	Y uint32 `json:"y,string"`
}

func (h *Handler) fibonacci(c *gin.Context) {
	var input FibonacciInput

	if err := c.BindJSON(&input); err != nil {
		c.Error(e.New(
			err,
			e.ErrUnprocessableEntity,
			http.StatusUnprocessableEntity,
			codes.InvalidArgument,
		))
		return
	}

	fibonacciSequence, err := h.Service.Calculate(input.X, input.Y)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, fibonacciSequence)
}
