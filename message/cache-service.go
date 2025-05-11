package message

import (
	"context"
	"fmt"
	"insider/redis"
	"time"
)

const (
	messageIdKey = "providerMessageId"
	sentOnKey    = "sentOn"
	ttl          = time.Hour
)

type redisDataService struct{}

func (r *redisDataService) Store(
	ctx context.Context,
	id int64,
	providerMessageId string,
	sentOn time.Time,
) error {
	redisInstance := redis.Instance()
	redisKey := r.key(id)
	cmd := redisInstance.HSet(ctx, redisKey, messageIdKey, providerMessageId)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	cmd = redisInstance.HSet(ctx, redisKey, sentOnKey, sentOn.Format(time.RFC3339))
	if cmd.Err() != nil {
		return cmd.Err()
	}

	bcmd := redisInstance.Expire(ctx, redisKey, ttl)
	if bcmd.Err() != nil {
		return bcmd.Err()
	}

	return nil
}

func (r *redisDataService) Exists(
	ctx context.Context,
	id int64,
) (*bool, error) {
	redisInstance := redis.Instance()
	redisKey := r.key(id)
	cmd := redisInstance.Exists(ctx, redisKey)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	result, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	var castedResult bool
	if result > 0 {
		castedResult = true
	} else {
		castedResult = false
	}

	return &castedResult, nil
}

func (r *redisDataService) key(id int64) string {
	return fmt.Sprintf("S:%d", id)
}
