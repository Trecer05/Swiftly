package chat

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"errors"

	"github.com/Trecer05/Swiftly/internal/config/logger"
	authErrors "github.com/Trecer05/Swiftly/internal/errors/auth"
	chatErrors "github.com/Trecer05/Swiftly/internal/errors/chat"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	kafkaModels "github.com/Trecer05/Swiftly/internal/model/kafka"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/chat"
	kafka "github.com/Trecer05/Swiftly/internal/repository/kafka/chat"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/chat"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"

	"github.com/gorilla/mux"
)

func GetTeamDashboardHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kfMgr *kafka.KafkaManager, ctx context.Context) {
	id, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID", authErrors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	user, err := mgr.GetStartUserInfo(id)
	if err != nil {
		logger.Logger.Error("Error getting user info", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	if len(user.Projects) != 0 {
	    for i := range user.Projects {
	        project := &user.Projects[i]
	        
	        if err := kfMgr.SendMessage(ctx, "dashboard", kafkaModels.TasksGet{UserID: user.ID, ProjectID: project.ID}); err != nil {
	            logger.Logger.Error("Error sending message", err)
	            continue
	        }
	        
	        data, err := kfMgr.WaitForResponse(user.ID, time.Second * 5)
	        if err != nil {
	            logger.Logger.Error("Error waiting for response", err)
	            serviceHttp.NewErrorBody(w, "application/json", err, http.StatusGatewayTimeout)
	            return
	        }
	
	        switch data.Type {
	        case "error":
	            var e kafkaModels.Error
	            json.Unmarshal(data.Payload, &e)
	            logger.Logger.Error("Error unmarshalling error response", e)
	        case "tasks":
	            var tasks []models.UserTask
	            json.Unmarshal(data.Payload, &tasks)
	            project.Tasks = append(project.Tasks, tasks...)
	        }
	    }
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func AddUserToTeamByUsernameHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, rds *redis.Manager) {
	vars := mux.Vars(r)
	
	teamID, err := strconv.Atoi(vars["team_id"])
	if err != nil {
		logger.Logger.Error("Error get team_id: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	username := vars["username"]
	
	ownerID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID", authErrors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	_, err = mgr.AddUserToTeamByUsername(teamID, ownerID, username)
	if err != nil {
		if err == chatErrors.ErrNoUser {
			logger.Logger.Error("Error not user found: ", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		}
		if err == chatErrors.ErrUserNotAOwner {
			logger.Logger.Error("Error user not a owner: ", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		
		logger.Logger.Error("Error add user to team: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func RemoveUserFromTeamByIDHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, rds *redis.Manager) {
	vars := mux.Vars(r)
	
	teamID, err := strconv.Atoi(vars["team_id"])
	if err != nil {
		logger.Logger.Error("Error get team_id: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		logger.Logger.Error("Error get user_id: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	ownerID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID", authErrors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	err = mgr.RemoveUserFromTeamByID(teamID, ownerID, userID)
	if err != nil {
		if err == chatErrors.ErrNoUser {
			logger.Logger.Error("Error not user found: ", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		}
		if err == chatErrors.ErrUserNotAOwner {
			logger.Logger.Error("Error user not a owner: ", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		
		logger.Logger.Error("Error remove user from team: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func UpdateUserRoleHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, rds *redis.Manager) {
	vars := mux.Vars(r)
	
	teamID, err := strconv.Atoi(vars["team_id"])
	if err != nil {
		logger.Logger.Error("Error get team_id: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		logger.Logger.Error("Error get user_id: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	ownerID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID", authErrors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	var role models.UserRoleEdit
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		logger.Logger.Error("Error decode role: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	err = mgr.UpdateUserRole(teamID, ownerID, userID, role.NewRole)
	if err != nil {
		if err == chatErrors.ErrNoUser {
			logger.Logger.Error("Error not user found: ", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		}
		if err == chatErrors.ErrUserNotAOwner {
			logger.Logger.Error("Error user not a owner: ", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		
		logger.Logger.Error("Error update user role: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func CreateTeamHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	var team models.TeamCreate
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		logger.Logger.Error("Error decode team: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	var ok bool
	team.OwnerID, ok = r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID", authErrors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	id, err := mgr.CreateTeam(&team)
	if err != nil {
		logger.Logger.Error("Error create team: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"team_id": id,
	})
}

func EditTeamHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, redis *redis.Manager) {
	var team models.TeamEdit
	err := json.NewDecoder(r.Body).Decode(&team)
	if err != nil {
		logger.Logger.Error("Error decode team: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	vars := mux.Vars(r)
	team.ID, err = strconv.Atoi(vars["team_id"])
	if err != nil {
		logger.Logger.Error("Error getting team ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	var ok bool
	team.OwnerID, ok = r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID", authErrors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	err = mgr.EditTeam(&team)
	if err != nil {
		if errors.Is(err, chatErrors.ErrProjectNotFound) {
			logger.Logger.Error("Error edit team: ", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		} else if errors.Is(err, chatErrors.ErrNoFieldsToUpdate) {
			logger.Logger.Info("No fields to update")
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNoContent)
			return
		} else if errors.Is(err, chatErrors.ErrUserNotAOwner) {
			logger.Logger.Error("Error edit team: ", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		
		logger.Logger.Error("Error edit team: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func DeleteTeamHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, kfMgr *kafka.KafkaManager) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["team_id"])
	if err != nil {
		logger.Logger.Error("Error getting team ID", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	var ok bool
	userID, ok := r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID", authErrors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	err = mgr.DeleteTeam(userID, id)
	if err != nil {
		if errors.Is(err, chatErrors.ErrProjectNotFound) {
			logger.Logger.Error("Error delete team: ", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		} else if errors.Is(err, chatErrors.ErrUserNotAOwner) {
			logger.Logger.Error("Error delete team: ", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		
		logger.Logger.Error("Error delete team: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	if err := kfMgr.SendMessage(ctx, "team_delete", kafkaModels.TeamTasksDelete{TeamID: id}); err != nil {
	    logger.Logger.Error("Error sending message", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
	    return
	}
	
	data, err := kfMgr.WaitForResponse(userID, time.Second * 5)
	if err != nil {
	    logger.Logger.Error("Error waiting for response", err)
	    serviceHttp.NewErrorBody(w, "application/json", err, http.StatusGatewayTimeout)
	    return
	}
	
	switch data.Type {
	case "error":
	    var e kafkaModels.Error
	    json.Unmarshal(data.Payload, &e)
	    logger.Logger.Error("Error unmarshalling error response", e)
	case "deleted":
	    var s kafkaModels.Success
	    json.Unmarshal(data.Payload, &s)
	    logger.Logger.Info(s.Msg)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func GetTeamInfoHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	id := mux.Vars(r)["team_id"]
	
	team, err := mgr.GetTeamInfo(id)
	if err != nil {
		if errors.Is(err, chatErrors.ErrProjectNotFound) {
			logger.Logger.Error("Error getting team info: ", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusNotFound)
			return
		} else if errors.Is(err, chatErrors.ErrUserNotAOwner) {
			logger.Logger.Error("Error getting team info: ", err)
			serviceHttp.NewErrorBody(w, "application/json", err, http.StatusForbidden)
			return
		}
		
		logger.Logger.Error("Error getting team info: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"team":   team,
	})
}
