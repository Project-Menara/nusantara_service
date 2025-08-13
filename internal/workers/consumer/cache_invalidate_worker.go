package consumer

import (
	"context"
	"encoding/json"
	"log"
	"nusantara_service/configs"
	"nusantara_service/internal/data/dataSources/rabbitmq"
	"nusantara_service/internal/workers/payload"
	"strings"

	"github.com/streadway/amqp"
)

func ConsumeCacheInvalidateQueue() {
	queue := rabbitmq.CacheInvalidateQueueName

	err := rabbitmq.ConsumeQueueAuto(queue, func(d amqp.Delivery) {
		var task payload.CacheInvalidateTask
		if err := json.Unmarshal(d.Body, &task); err != nil {
			log.Println("cache.invalidate unmarshal:", err)
			return
		}

		ctx := context.Background()
		rdb := configs.InitRedis()

		for _, key := range task.Keys {
			if hasWildcard(key) {
				iter := rdb.Scan(ctx, 0, key, 0).Iterator()
				for iter.Next(ctx) {
					_ = configs.DeleteRedis(ctx, iter.Val()).Error()
				}
				if err := iter.Err(); err != nil {
					log.Println("scan err:", err)
				}
			} else {
				_ = configs.DeleteRedis(ctx, key).Error()
			}
		}
	})
	if err != nil {
		log.Println("consume cache.invalidate.q err:", err)
	}
}

func hasWildcard(s string) bool {
	return strings.ContainsAny(s, "*?[")
}
