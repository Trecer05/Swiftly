package chat

import (
	"net/http"

	errors "github.com/Trecer05/Swiftly/internal/errors/auth"
	chatErrors "github.com/Trecer05/Swiftly/internal/errors/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
)

func ValidateGroupOwner(groupId int, r *http.Request, mgr *manager.Manager) (int, error) {
	userId, ok := r.Context().Value("id").(int)
	if !ok {
		return http.StatusUnauthorized, errors.ErrUnauthorized
	}

	ok, err := mgr.ValidateOwnerId(groupId, userId)
	switch {
	case err == chatErrors.ErrNoGroupFound:
		return http.StatusNotFound, err
	case err != nil:
		return http.StatusInternalServerError, err
	}

	if !ok {
		return http.StatusForbidden, errors.ErrGroupForbidden
	}
	return http.StatusOK, nil
}