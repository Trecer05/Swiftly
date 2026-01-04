package cloud

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"sync"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	"github.com/Trecer05/Swiftly/internal/filemanager/cloud"
	models "github.com/Trecer05/Swiftly/internal/model/kafka"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"

	"time"
)

type KafkaManager struct {
	Brokers []string
	Writer  *kafka.Writer
	Reader  *kafka.Reader

	mu        sync.Mutex
	responses map[string]chan models.KafkaResponse
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

	var reader *kafka.Reader
	for {
		reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:           brokers,
			Topic:             topic,
			GroupID:           groupID,
			MaxBytes:          10e6,
			MinBytes:          1,
			MaxWait:           3 * time.Second,
			ReadLagInterval:   -1,
			HeartbeatInterval: 3 * time.Second,
			CommitInterval:    0,
			RebalanceTimeout:  30 * time.Second,
			SessionTimeout:    30 * time.Second,
			StartOffset:       kafka.FirstOffset,
			Logger:            kafka.LoggerFunc(func(string, ...interface{}) {}),
			ErrorLogger:       kafka.LoggerFunc(logger.Logger.Errorf),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		msg, err := reader.FetchMessage(ctx)
		cancel()
		if err == nil || err == context.DeadlineExceeded || errors.Is(err, io.EOF) {
			if err == nil {
				reader.CommitMessages(ctx, msg)
			}
			break
		}
		logger.Logger.Warnf("Kafka not ready, retrying...: %v", err)
		reader.Close()
		time.Sleep(2 * time.Second)
	}

	return &KafkaManager{
		Brokers:   brokers,
		Writer:    writer,
		Reader:    reader,
		responses: make(map[string]chan models.KafkaResponse),
	}
}

func (km *KafkaManager) SendMessage(ctx context.Context, key string, value interface{}) (uuid.UUID, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return uuid.Nil, err
	}

	correlationID := uuid.New()
	header := kafka.Header{
		Key:   "correlation_id",
		Value: []byte(correlationID.String()),
	}

	msg := kafka.Message{
		Key:     []byte(key),
		Value:   []byte(data),
		Headers: []kafka.Header{header},
		Time:    time.Now(),
	}

	if err := km.Writer.WriteMessages(ctx, msg); err != nil {
		return uuid.Nil, err
	}

	logger.Logger.Infof("Sent message: key=%s value=%s with correlationID=%s", key, value, correlationID.String())
	return correlationID, nil
}

func (km *KafkaManager) ReadChatMessages(ctx context.Context) {
	for {
		msg, err := km.Reader.FetchMessage(ctx)

		if err != nil {
			if errors.Is(err, io.EOF) ||
				strings.Contains(err.Error(), "no messages") ||
				errors.Is(err, context.DeadlineExceeded) {
				continue
			}
			logger.Logger.Errorf("Kafka error: %v", err)
			continue
		}

		// var corrID string
		// for _, h := range msg.Headers {
		// 	if h.Key == "correlation_id" {
		// 		corrID = string(h.Value)
		// 		break
		// 	}
		// }

		switch string(msg.Key) {
		case "team_storage_create":
			var req models.TeamStorageCreate

			if err := json.Unmarshal(msg.Value, &req); err != nil {
				logger.Logger.Errorf("Error unmarshaling message: %v", err)

				data, _ := json.Marshal(models.Error{Err: err})

				km.SendMessage(ctx, "error", models.Envelope{Type: "error", Payload: data})
				continue
			}

			err := cloud.CreateTeamFolders(req.TeamID)
			if err != nil {
				logger.Logger.Errorf("Error create folders: %v", err)
				data, _ := json.Marshal(models.Error{Err: err})
				km.SendMessage(ctx, "error", models.Envelope{Type: "error", Payload: data})
				continue
			}

			km.SendMessage(ctx, "created", models.Envelope{Type: "tasks", Payload: nil})
		case "check_user_response":
			err := SendResponse(ctx, km, msg)
			if err != nil {
				logger.Logger.Errorf("Error sending envelope: %v", err)
				continue
			}
		}

	}
}

func SendResponse(ctx context.Context, km *KafkaManager, msg kafka.Message) error {
	var resp models.KafkaResponse
	if err := json.Unmarshal(msg.Value, &resp); err != nil {
		logger.Logger.Errorf("Error unmarshaling response message: %v", err)
		return err
	}

	km.mu.Lock()
	ch, ok := km.responses[resp.CorrelationID]
	km.mu.Unlock()

	if ok {
		ch <- resp
	} else {
		logger.Logger.Warnf("No waiting channel for correlationID=%s", resp.CorrelationID)
	}

	return nil
}

func (km *KafkaManager) WaitForResponse(correlationID string, timeout time.Duration) (models.KafkaResponse, error) {
	ch := make(chan models.KafkaResponse, 1)

	km.mu.Lock()
	km.responses[correlationID] = ch
	km.mu.Unlock()

	defer func() {
		km.mu.Lock()
		delete(km.responses, correlationID)
		km.mu.Unlock()
	}()

	select {
	case resp := <-ch:
		return resp, nil
	case <-time.After(timeout):
		return models.KafkaResponse{Err: context.DeadlineExceeded}, context.DeadlineExceeded
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
