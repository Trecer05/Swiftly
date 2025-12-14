package cloud

import (
	"context"
	"encoding/json"
	"sync"
	"io"
	"errors"
	"strings"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	models "github.com/Trecer05/Swiftly/internal/model/kafka"
	"github.com/Trecer05/Swiftly/internal/filemanager/cloud"
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
            Brokers: brokers,
            Topic: topic,
            GroupID: groupID,
            MaxBytes: 10e6,
            MinBytes: 1,
            MaxWait: 3 * time.Second,
            ReadLagInterval: -1,
            HeartbeatInterval: 3 * time.Second,
            CommitInterval: 0,
            RebalanceTimeout: 30 * time.Second,
            SessionTimeout: 30 * time.Second,
            StartOffset: kafka.FirstOffset,
            Logger: kafka.LoggerFunc(func(string, ...interface{}) {}),
            ErrorLogger: kafka.LoggerFunc(logger.Logger.Errorf),
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
        Brokers: brokers,
        Writer:  writer,
        Reader:  reader,
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

func (km *KafkaManager) ReadChatMessages(ctx context.Context) {
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
