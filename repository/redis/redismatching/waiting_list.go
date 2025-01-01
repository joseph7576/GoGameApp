package redismatching

import (
	"GoGameApp/entity"
	"GoGameApp/pkg/richerror"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	WaitingListPrefix = "waitinglist"
)

func (d DB) AddToWaitingList(userID uint, category entity.Category) error {
	const op = "redismatching.AddToWaitingList"

	err := d.adapter.Client().ZAdd(context.Background(),
		fmt.Sprintf("%s:%s", WaitingListPrefix, category),
		redis.Z{
			Score:  float64(time.Now().UnixMicro()),
			Member: fmt.Sprintf("%d", userID),
		}).Err()
	if err != nil {
		return richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}
