package consumer

import (
	"context"
	"encoding/json"
	"log"
	"nusantara_service/internal/data/dataSources/cloudinary"
	"nusantara_service/internal/data/dataSources/rabbitmq"
	"nusantara_service/internal/workers/payload"

	"github.com/streadway/amqp"
)

func ConsumeImageDeleteQueue() {
	queue := "image.delete.q"

	err := rabbitmq.ConsumeQueueAuto(queue, func(msg amqp.Delivery) {
		var task payload.ImageDeleteTask
		if err := json.Unmarshal(msg.Body, &task); err != nil {
			log.Println("image.delete.q unmarshal:", err)
			return
		}

		cld, err := cloudinary.NewCloudinaryService()
		if err != nil {
			log.Println("cloudinary init:", err)
			return
		}
		ctx := context.Background()

		type delRes struct {
			pid string
			err error
		}

		resCh := make(chan delRes)
		for _, pid := range task.PublicIDs {
			p := pid
			go func() {
				err := cld.DestroyImage(ctx, p)
				resCh <- delRes{pid: p, err: err}
			}()
		}

		for i := 0; i < len(task.PublicIDs); i++ {
			r := <-resCh
			if r.err != nil {
				log.Printf("delete fail publicID=%s err=%v\n", r.pid, r.err)
			} else {
				log.Printf("deleted publicID=%s\n", r.pid)
			}
		}
	})
	if err != nil {
		log.Println("consume image.delete.q err:", err)
	}
}
