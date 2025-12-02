package task_tracker

import (
	"net/http"
	"strconv"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/Trecer05/Swiftly/internal/config/logger"
	authErrors "github.com/Trecer05/Swiftly/internal/errors/auth"
	manager "github.com/Trecer05/Swiftly/internal/repository/postgres/task_tracker"
	redis "github.com/Trecer05/Swiftly/internal/repository/cache/task_tracker"
	serviceHttp "github.com/Trecer05/Swiftly/internal/transport/http"
	models "github.com/Trecer05/Swiftly/internal/model/task_tracker"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, rds *redis.WebSocketManager) {
	vars := mux.Vars(r)
	
	teamID, err := strconv.Atoi(vars["team_id"])
	if err != nil {
		logger.Logger.Error("Error getting team ID", authErrors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	columnID, err := strconv.Atoi(vars["column_id"])
	if err != nil {
		logger.Logger.Error("Error getting column ID", authErrors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	taskID, err := strconv.Atoi(vars["task_id"])
	if err != nil {
		logger.Logger.Error("Error getting task ID", authErrors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	columnPosition, err := mgr.DeleteTask(columnID, taskID)
	if err != nil {
		logger.Logger.Error("Error deleting task", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	msg := models.TaskDeleteMessage{
		ID: taskID,
		ColumnID: columnID,
		ColumnPosition: columnPosition,
		Message: "task deleted",
	}
	
	var envelope models.Envelope
	envelope.Data, err = json.Marshal(msg)
	if err != nil {
		logger.Logger.Error("Error marshaling task delete message", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	envelope.TeamID = teamID
	envelope.Type = models.TaskDeleteType
	
	rds.SendToUser(teamID, envelope)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, rds *redis.WebSocketManager) {
	var req models.CreateTaskRequest
	
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Logger.Error("Error decode task create request: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	var ok bool
	req.CreatorID, ok = r.Context().Value("id").(int)
	if !ok {
		logger.Logger.Error("Error getting user ID", authErrors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	vars := mux.Vars(r)
	req.ColumnID, err = strconv.Atoi(vars["column_id"])
	if err != nil {
		logger.Logger.Error("Error getting column ID", authErrors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	req.TeamID, err = strconv.Atoi(vars["team_id"])
	if err != nil {
		logger.Logger.Error("Error getting team ID", authErrors.ErrUnauthorized)
		serviceHttp.NewErrorBody(w, "application/json", authErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	
	taskID, err := mgr.CreateTask(&req)
	if err != nil {
		logger.Logger.Error("Error creating task: ", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	msg := models.TaskCreateMessage{
		ID: taskID,
		Title: req.Title,
		ColumnID: req.ColumnID,
		ColumnPosition: req.PositionInColumn,
		Deadline: req.Deadline,
		Message: "task created",
	}
	
	var envelope models.Envelope
	envelope.Data, err = json.Marshal(msg)
	if err != nil {
		logger.Logger.Error("Error marshaling task delete message", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	envelope.Type = models.TaskCreateType
	envelope.TeamID = req.TeamID
	
	rds.SendToUser(req.TeamID, envelope)
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": taskID,
	})
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	vars := mux.Vars(r)
	
	taskID, err := strconv.Atoi(vars["task_id"])
	if err != nil {
		logger.Logger.Error("Error converting task ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

    teamID, err := strconv.Atoi(vars["team_id"])
	if err != nil {
		logger.Logger.Error("Error converting team ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

    task, err := mgr.GetTaskByID(taskID, teamID)
	if err != nil {
		logger.Logger.Error("Error getting task", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task": task,
	})
}

func GetDashboardInfoHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager) {
	vars := mux.Vars(r)
	
	teamID, err := strconv.Atoi(vars["team_id"])
	if err != nil {
		logger.Logger.Error("Error converting team ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

    dashboardInfo, err := mgr.GetDashboardInfo(teamID)
	if err != nil {
		logger.Logger.Error("Error getting dashboard info", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"dashboard_info": dashboardInfo,
	})
}

func UpdateTaskTitleHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, rds *redis.WebSocketManager) {
	var req models.TaskTitleUpdateRequest
	var err error
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	vars := mux.Vars(r)
	
	req.ID, err = strconv.Atoi(vars["task_id"])
	if err != nil {
		logger.Logger.Error("Error converting task ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	teamID, err := strconv.Atoi(vars["team_id"])
	if err != nil {
		logger.Logger.Error("Error converting team ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	if req.ColumnID, req.PositionInColumn, err = mgr.UpdateTaskTitle(&req); err != nil {
		logger.Logger.Error("Error updating task", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	msg := models.TaskTitleUpdateMessage{
		ID: req.ID,
		Title: req.Title,
		ColumnID: req.ColumnID,
		ColumnPosition: req.PositionInColumn,
		Message: "task updated",
	}
	
	var envelope models.Envelope
	envelope.Data, err = json.Marshal(msg)
	if err != nil {
		logger.Logger.Error("Error marshaling task delete message", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	envelope.Type = models.TaskTitleUpdateType
	envelope.TeamID = teamID
	
	rds.SendToUser(teamID, envelope)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func UpdateTaskDeadlineHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, rds *redis.WebSocketManager) {
	var req models.TaskDeadlineUpdateRequest
	var err error
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	vars := mux.Vars(r)
	
	req.ID, err = strconv.Atoi(vars["task_id"])
	if err != nil {
		logger.Logger.Error("Error converting task ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	teamID, err := strconv.Atoi(vars["team_id"])
	if err != nil {
		logger.Logger.Error("Error converting team ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	if req.ColumnID, req.PositionInColumn, err = mgr.UpdateTaskDeadline(&req); err != nil {
		logger.Logger.Error("Error updating task", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	msg := models.TaskDeadlineUpdateMessage{
		ID: req.ID,
		Deadline: req.Deadline,
		ColumnID: req.ColumnID,
		ColumnPosition: req.PositionInColumn,
		Message: "task updated",
	}
	
	var envelope models.Envelope
	envelope.Data, err = json.Marshal(msg)
	if err != nil {
		logger.Logger.Error("Error marshaling task delete message", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	envelope.Type = models.TaskDeadlineUpdateType
	envelope.TeamID = teamID
	
	rds.SendToUser(teamID, envelope)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func UpdateTaskDeveloperHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, rds *redis.WebSocketManager) {
	var req models.TaskDeveloperUpdateRequest
	var err error
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	vars := mux.Vars(r)
	
	req.ID, err = strconv.Atoi(vars["task_id"])
	if err != nil {
		logger.Logger.Error("Error converting task ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	teamID, err := strconv.Atoi(vars["team_id"])
	if err != nil {
		logger.Logger.Error("Error converting team ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	if req.ColumnID, req.PositionInColumn, err = mgr.UpdateTaskDeveloper(&req); err != nil {
		logger.Logger.Error("Error updating task", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	msg := models.TaskDeveloperUpdateMessage{
		ID: req.ID,
		DeveloperID: req.DeveloperID,
		ColumnID: req.ColumnID,
		ColumnPosition: req.PositionInColumn,
		Message: "task updated",
	}
	
	var envelope models.Envelope
	envelope.Data, err = json.Marshal(msg)
	if err != nil {
		logger.Logger.Error("Error marshaling task delete message", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	envelope.Type = models.TaskDeveloperUpdateType
	envelope.TeamID = teamID
	
	rds.SendToUser(teamID, envelope)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func UpdateTaskStatusHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, rds *redis.WebSocketManager) {
	var req models.TaskStatusUpdateRequest
	var err error
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	vars := mux.Vars(r)
	
	req.ID, err = strconv.Atoi(vars["task_id"])
	if err != nil {
		logger.Logger.Error("Error converting task ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	teamID, err := strconv.Atoi(vars["team_id"])
	if err != nil {
		logger.Logger.Error("Error converting team ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	if req.ColumnID, req.PositionInColumn, err = mgr.UpdateTaskStatus(&req); err != nil {
		logger.Logger.Error("Error updating task", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	msg := models.TaskStatusUpdateMessage{
		ID: req.ID,
		Status: req.Status,
		ColumnID: req.ColumnID,
		ColumnPosition: req.PositionInColumn,
		Message: "task updated",
	}
	
	var envelope models.Envelope
	envelope.Data, err = json.Marshal(msg)
	if err != nil {
		logger.Logger.Error("Error marshaling task delete message", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	envelope.Type = models.TaskStatusUpdateType
	envelope.TeamID = teamID
	
	rds.SendToUser(teamID, envelope)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func UpdateTaskDescriptionHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, rds *redis.WebSocketManager) {
	var req models.TaskDescriptionUpdateRequest
	var err error
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	vars := mux.Vars(r)
	
	req.ID, err = strconv.Atoi(vars["task_id"])
	if err != nil {
		logger.Logger.Error("Error converting task ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	teamID, err := strconv.Atoi(vars["team_id"])
	if err != nil {
		logger.Logger.Error("Error converting team ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	if req.ColumnID, req.PositionInColumn, err = mgr.UpdateTaskDescription(&req); err != nil {
		logger.Logger.Error("Error updating task", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	msg := models.TaskDescriptionUpdateMessage{
		ID: req.ID,
		Description: req.Description,
		ColumnID: req.ColumnID,
		ColumnPosition: req.PositionInColumn,
		Message: "task updated",
	}
	
	var envelope models.Envelope
	envelope.Data, err = json.Marshal(msg)
	if err != nil {
		logger.Logger.Error("Error marshaling task delete message", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	envelope.Type = models.TaskDescriptionUpdateType
	envelope.TeamID = teamID
	
	rds.SendToUser(teamID, envelope)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func UpdateTaskColumnHandler(w http.ResponseWriter, r *http.Request, mgr *manager.Manager, rds *redis.WebSocketManager) {
	var req models.TaskColumnUpdateRequest
	var err error
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Logger.Error("Error decoding request body", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	vars := mux.Vars(r)
	
	req.ID, err = strconv.Atoi(vars["task_id"])
	if err != nil {
		logger.Logger.Error("Error converting task ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}
	
	teamID, err := strconv.Atoi(vars["team_id"])
	if err != nil {
		logger.Logger.Error("Error converting team ID to integer", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusBadRequest)
		return
	}

	if req.ColumnID, req.PositionInColumn, err = mgr.UpdateTaskColumn(&req); err != nil {
		logger.Logger.Error("Error updating task", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	msg := models.TaskColumnUpdateMessage{
		ID: req.ID,
		ColumnID: req.ColumnID,
		ColumnPosition: req.PositionInColumn,
		Message: "task updated",
	}
	
	var envelope models.Envelope
	envelope.Data, err = json.Marshal(msg)
	if err != nil {
		logger.Logger.Error("Error marshaling task delete message", err)
		serviceHttp.NewErrorBody(w, "application/json", err, http.StatusInternalServerError)
		return
	}
	
	envelope.Type = models.TaskColumnUpdateType
	envelope.TeamID = teamID
	
	rds.SendToUser(teamID, envelope)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}
