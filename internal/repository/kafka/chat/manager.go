package chat

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"sync"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	chatModels "github.com/Trecer05/Swiftly/internal/model/chat"
	models "github.com/Trecer05/Swiftly/internal/model/kafka"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	"github.com/segmentio/kafka-go"

	"time"
)

type KafkaManager struct {
	Brokers []string
	Writer  *kafka.Writer
	Reader  *kafka.Reader

	mu        sync.Mutex
	responses map[int]chan models.Envelope
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
		responses: make(map[int]chan models.Envelope),
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

func (km *KafkaManager) ReadMessage(ctx context.Context) error {
	msg, err := km.Reader.ReadMessage(ctx)
	if err != nil {
		return err
	}

	logger.Logger.Infof("Received message: key=%s value=%s", string(msg.Key), string(msg.Value))
	return nil
}

func (km *KafkaManager) ReadAuthMessages(ctx context.Context) {
	logger.Logger.Info("Starting to read messages")

	for {
		msg, err := km.Reader.FetchMessage(ctx)
		if errors.Is(err, io.EOF) ||
			strings.Contains(err.Error(), "no messages") ||
			errors.Is(err, context.DeadlineExceeded) {
			continue
		}
		if err != nil {
			logger.Logger.Errorf("Kafka error: %v", err)
			continue
		}

		switch string(msg.Key) {
		case "password":
			err := SendEnvel(ctx, km, "password", msg)
			if err != nil {
				logger.Logger.Errorf("Error sending envelope: %v", err)
				continue
			}
		case "phone":
			err := SendEnvel(ctx, km, "phone", msg)
			if err != nil {
				logger.Logger.Errorf("Error sending envelope: %v", err)
				continue
			}
		case "email":
			err := SendEnvel(ctx, km, "email", msg)
			if err != nil {
				logger.Logger.Errorf("Error sending envelope: %v", err)
				continue
			}
		}
	}
}

func (km *KafkaManager) ReadTaskMessages(ctx context.Context) {
	logger.Logger.Info("Starting to read messages")

	for {
		msg, err := km.Reader.FetchMessage(ctx)
		if errors.Is(err, io.EOF) ||
			strings.Contains(err.Error(), "no messages") ||
			errors.Is(err, context.DeadlineExceeded) {
			continue
		}
		if err != nil {
			logger.Logger.Errorf("Kafka error: %v", err)
			continue
		}

		switch string(msg.Key) {
		case "tasks":
			err := SendEnvel(ctx, km, "tasks", msg)
			if err != nil {
				logger.Logger.Errorf("Error sending envelope: %v", err)
				continue
			}
		case "error":
			err := SendEnvel(ctx, km, "error", msg)
			if err != nil {
				logger.Logger.Errorf("Error sending envelope: %v", err)
				continue
			}
		}
	}
}

func (km *KafkaManager) ReadCloudMessages(ctx context.Context, manager *manager.Manager) {
	logger.Logger.Info("Starting to read messages")

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

		switch string(msg.Key) {
		case "created":
			err := SendEnvel(ctx, km, "tasks", msg)
			if err != nil {
				logger.Logger.Errorf("Error sending envelope: %v", err)
				continue
			}
		case "check_user":
			logger.Logger.Info("Recieved check_user message")
			var req models.CheckUserInTeam

			// извлекаем данные из msg.Value
			if err := json.Unmarshal(msg.Value, &req); err != nil {
				logger.Logger.Errorf("Error unmarshaling message: %v", err)

				data, _ := json.Marshal(models.Error{Err: err})

				km.SendMessage(ctx, "error", models.Envelope{Type: "error", Payload: data})
				continue
			}

			// извлекаем correlation_id из заголовков (если есть) чтобы вернуть ответ с тем же id
			var corrID string
			for _, h := range msg.Headers {
				if h.Key == "correlation_id" {
					corrID = string(h.Value)
					break
				}
			}

			logger.Logger.Info("Correlation ID: " + corrID)

			isInTeam, err := manager.IsUserInTeam(req.TeamID, req.UserID)
			if err != nil {
				logger.Logger.Errorf("Error checking user in team: %v", err)
				data, _ := json.Marshal(models.Error{Err: err})
				km.SendMessage(ctx, "error", models.Envelope{Type: "error", Payload: data})
				continue
			}
			respPayload, err := json.Marshal(chatModels.CheckUserInTeamResponse{
				IsInTeam: isInTeam,
			})

			if err != nil {
				logger.Logger.Errorf("Error marshaling response: %v", err)
				// data, _ := json.Marshal(models.Error{Err: err})
				km.SendMessage(ctx, "error", models.KafkaResponse{CorrelationID: corrID, Err: err})
				continue
			}

			headers := []kafka.Header{}
			if corrID != "" {
				headers = append(headers, kafka.Header{Key: "correlation_id", Value: []byte(corrID)})
			}

			// reply := kafka.Message{
			// 	Key:     []byte("check_user_response"),
			// 	Value:   respPayload,
			// 	Headers: headers,
			// 	Time:    time.Now(),
			// }
			km.SendMessage(ctx, "check_user_response", models.KafkaResponse{
				CorrelationID: corrID,
				Payload:       respPayload,
				Err:           nil,
			})
			// if err := km.Writer.WriteMessages(ctx, reply); err != nil {
			// 	logger.Logger.Errorf("error writing check_user response: %v", err)
			// }
			km.Reader.CommitMessages(ctx, msg)
		case "error":
			err := SendEnvel(ctx, km, "error", msg)
			if err != nil {
				logger.Logger.Errorf("Error sending envelope: %v", err)
				continue
			}
		}
	}
}

func SendEnvel(ctx context.Context, km *KafkaManager, key string, msg kafka.Message) error {
	var envel models.Envelope
	if err := json.Unmarshal(msg.Value, &envel); err != nil {
		logger.Logger.Errorf("Error unmarshaling message: %v", err)
		return err
	}

	userID := *envel.UserID

	km.mu.Lock()
	ch, ok := km.responses[userID]
	km.mu.Unlock()

	if ok {
		ch <- envel
	} else {
		logger.Logger.Warnf("No waiting channel for userID=%d", userID)
	}

	return nil
}

func (km *KafkaManager) WaitForResponse(userID int, timeout time.Duration) (models.Envelope, error) {
	ch := make(chan models.Envelope, 1)

	km.mu.Lock()
	km.responses[userID] = ch
	km.mu.Unlock()

	defer func() {
		km.mu.Lock()
		delete(km.responses, userID)
		km.mu.Unlock()
	}()

	select {
	case resp := <-ch:
		return resp, nil
	case <-time.After(timeout):
		return models.Envelope{Type: "error"}, context.DeadlineExceeded
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
