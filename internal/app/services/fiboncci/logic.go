package fiboncci

import (
	"errors"
	"fmt"
	"github.com/VSKrivoshein/FBS-test/internal/app/cache"
	e "github.com/VSKrivoshein/FBS-test/internal/app/custom_err"
	"google.golang.org/grpc/codes"
	"math/big"
	"net/http"
)

const (
	MaxIntFibonacci = uint32(5000)
)

type service struct {
	cache cache.Rdb
}

func NewService(rdb cache.Rdb) Service {
	return &service{
		cache: rdb,
	}
}

func (s *service) Calculate(startInd uint32, endInd uint32) ([]string, error) {
	// Валидация значений
	if startInd > endInd {
		return nil, e.New(
			errors.New("startInd > endInd"),
			e.ErrSecondNumberLessThanFirst,
			http.StatusUnprocessableEntity,
			codes.InvalidArgument,
		)
	}
	// Сообщение grpc ограничено по размеру
	if endInd >= MaxIntFibonacci {
		return nil, e.New(
			errors.New("endInd > MaxIntFibonacci"),
			e.ErrMaxIntFibonacci,
			http.StatusUnprocessableEntity,
			codes.InvalidArgument,
		)
	}

	fibonacciSequence, lastCachedInd, err := s.cache.GetFibonacciSequence(startInd, endInd)
	if err != nil {
		return nil, fmt.Errorf(e.GetInfo(), err)
	}

	// Все необходимые данные были в кэше, сразу отдаем результат
	if lastCachedInd == endInd {
		return fibonacciSequence, nil
	}

	// Получаем 2 последних известных числа для расчетов
	first, ok := new(big.Int).SetString(fibonacciSequence[len(fibonacciSequence)-2], 10)
	second, ok := new(big.Int).SetString(fibonacciSequence[len(fibonacciSequence)-1], 10)
	if !ok {
		return nil, e.New(
			errors.New("error during creating new big int"),
			e.ErrUnexpected,
			http.StatusInternalServerError,
			codes.Internal,
		)
	}

	// нужно дождаться пока кэширование закончится чтобы не потерять данные из fibonacciSequence
	// при этом необходимо отдать fibonacciSequence сразу как завершаться расчеты
	// поэтому возвращаем результат и дожидаемся окончания кэширования
	done := make(chan struct{})
	defer func() {
		<-done
	}()

	// Запускаем goroutine, которая проверяет появление новых значений в slice и записывает их в redis
	go s.cache.SetFibonacciSequence(&fibonacciSequence, startInd, endInd, done)

	// Рассчитываем Фибоначи
	for i := lastCachedInd; i <= endInd; i++ {
		currentNumber := new(big.Int)
		first, second = second, currentNumber.Add(first, second)
		fibonacciSequence = append(fibonacciSequence, second.String())
	}

	// Если у нас массив больше необходимого, то возвращаем только нужное
	if lastCachedInd <= startInd {
		return fibonacciSequence[(startInd - lastCachedInd + 2):], nil
	}

	return fibonacciSequence, nil
}
