package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	models "github.com/Trecer05/Swiftly/internal/model/kafka"
	postgres "github.com/Trecer05/Swiftly/internal/repository/postgres/auth"

	"github.com/segmentio/kafka-go"
)

type KafkaManager struct {
    Brokers []string
    Writer  *kafka.Writer
    Reader  *kafka.Reader
}

func NewKafkaManager(brokers []string, topic, groupID string) *KafkaManager {
	logger.Logger.Infof("Create new KafkaManager with brokers: %v, topic: %s, groupID: %s", brokers, topic, groupID)

    writer := &kafka.Writer{
        Addr:         kafka.TCP(brokers...),
        Topic:        topic,
        Balancer:     &kafka.LeastBytes{},
        RequiredAcks: kafka.RequireAll,
        Async:        false,
    }

    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:  brokers,
        Topic:    topic,
        GroupID:  groupID,
        MinBytes: 10e3,
        MaxBytes: 10e6,
    })

    return &KafkaManager{
        Brokers: brokers,
        Writer:  writer,
        Reader:  reader,
    }
}

func (km *KafkaManager) SendMessage(ctx context.Context, key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }

    msg := kafka.Message{
        Key:   []byte(key),
        Value: []byte(data),
        Time:  time.Now(),
    }

    if err := km.Writer.WriteMessages(ctx, msg); err != nil {
        return err
    }

    logger.Logger.Infof("Sent message: key=%s value=%s", key, value)
    return nil
}

func (km *KafkaManager) ReadMessage(ctx context.Context) (*kafka.Message, error) {
    msg, err := km.Reader.ReadMessage(ctx)
    if err != nil {
        return nil, err
    }

    logger.Logger.Infof("Received message: key=%s value=%s", string(msg.Key), string(msg.Value))
    return &msg, err
}

func (km *KafkaManager) ReadUserEditMessages(ctx context.Context, mgr *postgres.Manager) {
	logger.Logger.Info("Starting to read messages")

	for {
		msg, err := km.Reader.ReadMessage(ctx)
		if err != nil {
			logger.Logger.Errorf("Error reading message: %v", err)
			continue
		}

		switch string(msg.Key) {
        case "password":
            var req models.PasswordEdit

            if err := json.Unmarshal(msg.Value, &req); err != nil {
                logger.Logger.Errorf("Error unmarshaling message: %v", err)

                data, _ := json.Marshal(models.Error{Err: err})

                km.SendMessage(ctx, "password", models.Envelope{Type: "error", Payload: data})
                continue
            }

            go km.PasswordEdit(ctx, mgr, req)
        case "phone":
            var req models.PhoneEdit

            if err := json.Unmarshal(msg.Value, &req); err != nil {
                logger.Logger.Errorf("Error unmarshaling message: %v", err)

                data, _ := json.Marshal(models.Error{Err: err})

                km.SendMessage(ctx, "phone", models.Envelope{Type: "error", Payload: data})
                continue
            }

            mgr.EditPhone(&req)
        case "email":
            var req models.EmailEdit

            if err := json.Unmarshal(msg.Value, &req); err != nil {
                logger.Logger.Errorf("Error unmarshaling message: %v", err)

                data, _ := json.Marshal(models.Error{Err: err})

                km.SendMessage(ctx, "email", models.Envelope{Type: "error", Payload: data})
                continue
            }

            mgr.EditEmail(&req)
        }
	}
}

func (km *KafkaManager) Close() error {
    if err := km.Writer.Close(); err != nil {
        return err
    }
    if err := km.Reader.Close(); err != nil {
        return err
    }
    return nil
}
