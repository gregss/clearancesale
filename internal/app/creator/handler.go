package creator

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"sale/internal/app"
)

type CreateHandler struct {
	App app.App
}

func (h *CreateHandler) Handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)

		var message app.CreateMessage
		if err := json.Unmarshal(d.Body, &message); err != nil {
			// logger.error()
			return
		}

		h.App.FillProductsByRegions(message)

		err := d.Ack(false)
		if err != nil {
			done <- err
		}
	}

	done <- nil
}
