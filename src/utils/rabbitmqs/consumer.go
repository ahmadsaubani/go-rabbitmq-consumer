package rabbitmqs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RPCHandler func(requestBody []byte) ([]byte, error)

func StartRPCConsumer(queueName, consumerName string, handler RPCHandler) error {
	ch, err := GetChannel()
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// QoS
	err = ch.Qos(
		1, // prefetchCount, hanya satu pesan yang akan dikirim ke consumer pada satu waktu
		0, //prefetchSize tidak ada batas ukuran byte dari pesan yang dikirim ke consumer
		false)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := ch.Consume(
		queueName,
		consumerName,
		false, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	// kode dijalankan paralel agar tidak menghalangi proses utama
	go func() {
		log.Printf("[*] RPC Consumer started on queue '%s'", queueName)
		for d := range msgs {
			log.Printf("[x] Received RPC request on %s: %s", queueName, string(d.Body))

			// memberikan timeout agar consumer stop merespons kalau RabbitMQ terhenti.
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

			responseBody, err := handler(d.Body)
			if err != nil {
				responseBody, _ = json.Marshal(map[string]interface{}{
					"success": false,
					"error":   err.Error(),
				})
			}

			if d.ReplyTo != "" {
				err = ch.PublishWithContext(ctx,
					"",
					d.ReplyTo,
					false,
					false,
					amqp.Publishing{
						ContentType:   "application/json",
						CorrelationId: d.CorrelationId,
						Body:          responseBody,
					})

				if err != nil {
					log.Printf("[!] Failed to reply to RPC request: %v", err)
				}
			}

			d.Ack(false)
			cancel()
		}
	}()

	return nil
}
