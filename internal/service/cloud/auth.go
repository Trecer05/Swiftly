package cloud

import (
	"net/http"

	errorAuthTypes "github.com/Trecer05/Swiftly/internal/errors/auth"
	kafkaManager "github.com/Trecer05/Swiftly/internal/repository/kafka/cloud"
)

func AuthorizeUserInTeam(r *http.Request, teamID int, kafkaManager *kafkaManager.KafkaManager) (int, error) {
	userID, ok := r.Context().Value("id").(int)
	if !ok {
		return 0, errorAuthTypes.ErrUnauthorized
	}
	if err := CheckUserInTeam(teamID, userID, kafkaManager); err != nil {
		return 0, err
	}
	return userID, nil
}
