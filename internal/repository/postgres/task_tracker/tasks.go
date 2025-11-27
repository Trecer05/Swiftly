package task_tracker

import (
	"database/sql"

	models "github.com/Trecer05/Swiftly/internal/model/chat"
	taskModels "github.com/Trecer05/Swiftly/internal/model/task_tracker"
)

func (manager *Manager) GetUserTasks(userID int, projectID int) ([]models.UserTask, error) {
    rows, err := manager.Conn.Query(`
        SELECT 
            t.id,
            t.title,
            t.description,
            t.label,
            t.completed_at,
            pd.level,
            pd.title as priority_title,
            pd.color as priority_color,
            t.status
        FROM tasks t
        INNER JOIN project_tasks pt ON t.id = pt.task_id
        INNER JOIN priority_definitions pd ON t.priority = pd.level
        WHERE pt.project_id = $1 
          AND (t.author_id = $2 OR t.developer_id = $2)
          AND t.status = 'in_progress'
        ORDER BY t.created_at DESC`, 
        projectID, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []models.UserTask

    for rows.Next() {
        var task models.UserTask
        var completedAt sql.NullTime
        var priorityLevel string
        var status string

        err := rows.Scan(
            &task.ID,
            &task.Name,
            &task.Description,
            &task.Label,
            &completedAt,
            &priorityLevel,
            &task.Priority.Title,
            &task.Priority.HexColor,
            &status,
        )
        if err != nil {
            return nil, err
        }

        switch priorityLevel {
        case "low":
            task.Priority.Type = models.PriorityLow
        case "medium":
            task.Priority.Type = models.PriorityMedium
        case "high":
            task.Priority.Type = models.PriorityHigh
        }

        if completedAt.Valid {
            task.EndTime = completedAt.Time
        }

        tasks = append(tasks, task)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return tasks, nil
}

func (manager *Manager) DeleteTeamTasks(teamID int) error {
    _, err := manager.Conn.Exec("DELETE FROM project_columns WHERE project_id = $1", teamID)
    return err
}

func (manager *Manager) CreateStartTasksTables(userID, projectID int) error {
    _, err := manager.Conn.Exec(`
    	INSERT INTO project_columns (project_id, title, position, created_by)
    	VALUES ($1, 'Start', 'Start tasks', $2)
    `, projectID, userID)
    if err != nil {
        return err
    }

    return nil
}

func (manager *Manager) GetTaskByID(taskID int) (taskModels.Task, error) {
	var task taskModels.Task
	var completedAt sql.NullTime
	var priorityLevel string
	var status string

	err := manager.Conn.QueryRow("SELECT id, author_id, developer_id, title, description, label, completed_at, priority_level, priority_title, priority_hex_color, status FROM tasks WHERE id = $1", taskID).Scan(
		&task.ID,
		&task.AuthorID,
		&task.DeveloperID,
		&task.Title,
		&task.Description,
		&task.Label,
		&completedAt,
		&priorityLevel,
		&task.Priority.Title,
		&task.Priority.HexColor,
		&status,
	)

	if err != nil {
		return task, err
	}

	switch priorityLevel {
	case "low":
		task.Priority.Type = taskModels.PriorityLow
	case "medium":
		task.Priority.Type = taskModels.PriorityMedium
	case "high":
		task.Priority.Type = taskModels.PriorityHigh
	}

	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}

	task.Status = status

	return task, nil
}
