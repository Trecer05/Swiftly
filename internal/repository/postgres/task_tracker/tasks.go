package task_tracker

import (
	"database/sql"

	models "github.com/Trecer05/Swiftly/internal/model/chat"
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
