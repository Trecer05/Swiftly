package auth

import (
	"context"
	"encoding/json"

	models "github.com/Trecer05/Swiftly/internal/model/kafka"
	postgres "github.com/Trecer05/Swiftly/internal/repository/postgres/auth"
)

func (km *KafkaManager) ValidatePasswordAndEdit(mgr *postgres.Manager, req models.PasswordEdit) error {
	err := mgr.ValidatePasswordByID(req.UserID, req.OldPassword)
	if err != nil {
		return err
	}

	err = mgr.EditPassword(&req)
	if err != nil {
		return err
	}

	return nil
}

func (km *KafkaManager) PasswordEdit(ctx context.Context, mgr *postgres.Manager, req models.PasswordEdit) {
	if err := km.ValidatePasswordAndEdit(mgr, req); err != nil {
		data, _ := json.Marshal(models.Error{Err: err})

		km.SendMessage(ctx, "password", models.Envelope{Type: "error", Payload: data})
		return
	}

	data, _ := json.Marshal(models.Status{Status: "ok"})

	if err := km.SendMessage(ctx, "password", models.Envelope{Type: "status", Payload: data}); err != nil {
		data, _ := json.Marshal(models.Error{Err: err})

		km.SendMessage(ctx, "password", models.Envelope{Type: "error", Payload: data})
		return
	}
}
