package cache

import (
	context "context"
	"errors"
	"fmt"
	e "github.com/VSKrivoshein/FBS-test/internal/app/custom_err"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	cacheKey = "cache"
)

type rdb struct {
	sync.RWMutex
	offset uint32
	client *redis.Client
}

func NewRdb() (Rdb, error) {
	rdb := &rdb{
		offset: 0,
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%v:%v",os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		}),
	}
	log.Info(fmt.Sprintf("%v:%v",os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")))

	ctx := context.Background()
	_, err := rdb.client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis connection ping failed: %w", err)
	}

	return rdb, nil
}

func (r *rdb) PrepareRDB() error {
	ctx := context.Background()
	length, err := r.client.LLen(ctx, cacheKey).Result()
	if err != nil {
		return e.New(
			err,
			errors.New("redis connection failed"),
			http.StatusInternalServerError,
			codes.Internal,
		)
	}
	r.offset = uint32(length)

	if r.offset < 2 {
		// Выставляем первые значения последовательности значения
		if err := r.client.RPush(ctx, cacheKey, "0", "1").Err(); err != nil {
			return e.New(
				err,
				errors.New("unable set to redis default values"),
				http.StatusInternalServerError,
				codes.Internal,
			)
		}
		r.offset = 2
	}

	return nil
}

func (r *rdb) GetFibonacciSequence(x uint32, y uint32) ([]string, uint32, error) {
	// Для использования RLock фиксируем значение offset
	stableOffset := r.offset
	var firstIndex uint32
	var res []string
	var err error
	ctx := context.Background()

	// Определяем какой массив возможно запросить в кэше, чтобы не тащить весь сразу, если нам нужна только часть
	if x <= stableOffset && y <= stableOffset {
		res, err = r.client.LRange(ctx, cacheKey, int64(x), int64(y)).Result()
		firstIndex = y
	} else if x <= stableOffset-2 && y > stableOffset {
		res, err = r.client.LRange(ctx, cacheKey, int64(x), int64(stableOffset)).Result()
		firstIndex = stableOffset
	} else {
		res, err = r.client.LRange(ctx, cacheKey, int64(stableOffset-2), int64(stableOffset)).Result()
		firstIndex = stableOffset
	}

	if err != nil {
		return nil, 0, e.New(
			err,
			e.ErrUnexpected,
			http.StatusInternalServerError,
			codes.Internal,
		)
	}

	if len(res) == 0 {
		//не найден кэш который должен там быть
		//заново получаем данные с редиса по кэшу для синхронизации
		if err := r.PrepareRDB(); err != nil {
			return nil, 0, fmt.Errorf(e.GetInfo(), err)
		}
		//Повторно делаем запрос в редис
		return r.GetFibonacciSequence(x, y)
	}

	return res, firstIndex, nil
}

func (r *rdb) SetFibonacciSequence(sequence *[]string, x uint32, y uint32, done chan struct{}) {
	// Проверяем, есть ли новые данные для кэша
	if y <= r.offset {
		return
	}

	ctx := context.Background()

	// Блокируем указатель на последний элемент в кэше (offset),
	// до тех пор пока не запишем все данные в кэш
	r.RLock()
	defer r.RUnlock()

	// Slice для кэша может отличаться по размерам указателей x/y, определяем верный индекс
	var indForCache uint32
	if x+2 < r.offset {
		indForCache = 0
	} else {
		indForCache = 2
	}

	// Запускаем цикл, для постоянной проверки на наличие новых значений в слайсе для кэша
	for {
		// Пишем в кэш новые значения (параллельно в массив пишутся новые значения)
		for ; indForCache < uint32(len(*sequence)); indForCache++ {
			r.client.RPush(ctx, cacheKey, (*sequence)[indForCache])
			r.offset++
		}

		// Выходим если больше не ожидаем появление новых значений
		if y <= r.offset {
			break
		}

		time.Sleep(time.Millisecond)
	}

	done <- struct{}{}
}
