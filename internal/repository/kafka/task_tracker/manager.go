package task_tracker

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	models "github.com/Trecer05/Swiftly/internal/model/kafka"
	mgr "github.com/Trecer05/Swiftly/internal/repository/postgres/task_tracker"
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

func (km *KafkaManager) ReadChatMessages(ctx context.Context, manager *mgr.Manager) {
    for {
        msg, err := km.Reader.ReadMessage(ctx)
        if err != nil {
            logger.Logger.Errorf("Error reading message: %v", err)
            return
        }

        switch string(msg.Key) {
        case "dashboard":
        	var req models.TasksGet
         
        	if err := json.Unmarshal(msg.Value, &req); err != nil {
			    logger.Logger.Errorf("Error unmarshaling message: %v", err)
			
			    data, _ := json.Marshal(models.Error{Err: err})
			
			    km.SendMessage(ctx, "error", models.Envelope{Type: "error", Payload: data})
			    continue
			}
			
			tasks, err := manager.GetUserTasks(req.UserID, req.ProjectID)
			if err != nil {
				logger.Logger.Errorf("Error getting tasks: %v", err)
				data, _ := json.Marshal(models.Error{Err: err})
				km.SendMessage(ctx, "error", models.Envelope{Type: "error", Payload: data})
				continue
			}
			
			data, _ := json.Marshal(tasks)
			km.SendMessage(ctx, "tasks", models.Envelope{Type: "tasks", Payload: data})
		case "team_delete":
		    var req models.TeamTasksDelete
		    
		    if err := json.Unmarshal(msg.Value, &req); err != nil {
			    logger.Logger.Errorf("Error unmarshaling message: %v", err)
			
			    data, _ := json.Marshal(models.Error{Err: err})
			
			    km.SendMessage(ctx, "error", models.Envelope{Type: "error", Payload: data})
			    continue
			}
			
			err := manager.DeleteTeamTasks(req.TeamID)
			if err != nil {
				logger.Logger.Errorf("Error deleting team: %v", err)
				data, _ := json.Marshal(models.Error{Err: err})
				km.SendMessage(ctx, "error", models.Envelope{Type: "error", Payload: data})
				continue
			}
			
			data, _ := json.Marshal(models.Success{Msg: "Team deleted successfully"})
			km.SendMessage(ctx, "success", models.Envelope{Type: "success", Payload: data})
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
