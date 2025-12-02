package task_tracker

import (
	"database/sql"
	"fmt"
	"time"
	"sort"

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

func (manager *Manager) GetTaskByID(taskID, teamID int) (taskModels.Task, error) {
	var task taskModels.Task
	var completedAt sql.NullTime
	var priorityLevel string
	var status string

	err := manager.Conn.QueryRow("SELECT id, author_id, developer_id, title, description, label, completed_at, priority_level, priority_title, priority_hex_color, status, deadline FROM tasks WHERE id = $1", taskID).Scan(
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
		&task.Deadline,
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

func (manager *Manager) CreateTask(req *taskModels.CreateTaskRequest) (int, error) {
    tx, err := manager.Conn.Begin()
    if err != nil {
        return 0, fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    var columnExists bool
    err = tx.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM project_columns pc 
            WHERE pc.id = $1 AND pc.project_id = $2
        )`, req.ColumnID, req.TeamID).Scan(&columnExists)
    if err != nil {
        return 0, fmt.Errorf("failed to check column: %w", err)
    }
    if !columnExists {
        return 0, fmt.Errorf("column not found or doesn't belong to team")
    }

    // TODO: нужно как то разработать систему получения пользователя из команды (возможно дубляж таблиц, но хз)
    var developerID *int
    // if req.ExecutorUsername != "" {
    //     err = tx.QueryRow(`
    //         SELECT id FROM users WHERE username = $1
    //     `, req.ExecutorUsername).Scan(&developerID)
    //     if err != nil {
    //         if err == sql.ErrNoRows {
    //             return 0, fmt.Errorf("executor user not found: %s", req.ExecutorUsername)
    //         }
    //         return 0, fmt.Errorf("failed to find executor: %w", err)
    //     }
    // }

    position := req.PositionInColumn
    if position == 0 {
        err = tx.QueryRow(`
            SELECT COALESCE(MAX(position_in_column), 0) + 1 
            FROM tasks WHERE column_id = $1
        `, req.ColumnID).Scan(&position)
        if err != nil {
            return 0, fmt.Errorf("failed to get next position: %w", err)
        }
    } else {
        _, err = tx.Exec(`
            UPDATE tasks 
            SET position_in_column = position_in_column + 1 
            WHERE column_id = $1 AND position_in_column >= $2
        `, req.ColumnID, position)
        if err != nil {
            return 0, fmt.Errorf("failed to shift tasks positions: %w", err)
        }
    }

    var taskID int
    query := `
        INSERT INTO tasks (
            author_id, developer_id, column_id, position_in_column, 
            title, description, label, deadline
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id
    `
    
    err = tx.QueryRow(query,
        req.CreatorID,
        developerID,
        req.ColumnID,
        position,
        req.Title,
        req.Description,
        req.Label,
        req.Deadline,
    ).Scan(&taskID)
    if err != nil {
        return 0, fmt.Errorf("failed to create task: %w", err)
    }

    _, err = tx.Exec(`
        INSERT INTO project_tasks (project_id, task_id) 
        VALUES ($1, $2)
    `, req.TeamID, taskID)
    if err != nil {
        return 0, fmt.Errorf("failed to link task to project: %w", err)
    }

    if err = tx.Commit(); err != nil {
        return 0, fmt.Errorf("failed to commit transaction: %w", err)
    }

    return taskID, nil
}

func (manager *Manager) DeleteTask(columnID, taskID int) (int, error) {
    var position int
    err := manager.Conn.QueryRow(`
        DELETE FROM tasks WHERE id = $1 AND column_id = $2 RETURNING position_in_column
    `, taskID, columnID).Scan(&position)
    if err != nil {
        return 0, fmt.Errorf("failed to delete task: %w", err)
    }
    
    return position, nil
}

func (manager *Manager) GetDashboardInfo(teamID int) (taskModels.DashboardInfo, error) {
    var info taskModels.DashboardInfo
    info.TeamID = teamID

    query := `
        SELECT 
            c.id AS column_id,
            c.title AS column_title,
            c.position AS column_position,
            c.color AS column_color,

            t.id AS task_id,
            t.title AS task_title,
            t.description AS task_description,
            t.status AS task_status,
            t.priority AS task_priority,
            t.label AS task_label,
            t.deadline AS task_deadline

        FROM project_columns c
        LEFT JOIN tasks t ON t.column_id = c.id
        WHERE c.project_id = $1
        ORDER BY c.position, t.position_in_column;
    `

    rows, err := manager.Conn.Query(query, teamID)
    if err != nil {
        return info, err
    }
    defer rows.Close()

    columnsMap := make(map[int]*taskModels.Column)

    for rows.Next() {
        var (
            colID       int
            colTitle    string
            colPos      int
            colColor    string

            taskPriority taskModels.Priority
            taskID      sql.NullInt64
            taskTitle   sql.NullString
            taskDesc    sql.NullString
            taskStatus  sql.NullString
            taskLabel   sql.NullString
            taskDeadline sql.NullTime
        )

        err := rows.Scan(
            &colID, &colTitle, &colPos, &colColor,
            &taskID, &taskTitle, &taskDesc, &taskStatus, &taskPriority.Type, &taskLabel, &taskDeadline,
        )
        if err != nil {
            return info, err
        }

        col, exists := columnsMap[colID]
        if !exists {
            col = &taskModels.Column{
                ID:       colID,
                Title:    colTitle,
                Position: colPos,
                HexColor: colColor,
                Tasks:    []taskModels.DashboardTask{},
            }
            columnsMap[colID] = col
        }

        if !taskID.Valid {
            continue
        }

        var deadline time.Time
        if taskDeadline.Valid {
            deadline = taskDeadline.Time
        }

        task := taskModels.DashboardTask{
            ID:          int(taskID.Int64),
            Title:       taskTitle.String,
            Description: nil,
            Status:      taskStatus.String,
            Priority:    taskPriority,
            Label:       taskLabel.String,
            Deadline:    deadline,
        }

        if taskDesc.Valid {
            d := taskDesc.String
            task.Description = &d
        }

        col.Tasks = append(col.Tasks, task)
    }

    for _, c := range columnsMap {
        info.Columns = append(info.Columns, *c)
    }

    sort.Slice(info.Columns, func(i, j int) bool {
        return info.Columns[i].Position < info.Columns[j].Position
    })

    return info, nil
}

func (manager *Manager) UpdateTaskTitle(req *taskModels.TaskTitleUpdateRequest) (int, int, error) {
	err := manager.Conn.QueryRow("UPDATE tasks SET title = $1 WHERE id = $2 RETURNING column_id, position_in_column", req.Title, req.ID).Scan(&req.ColumnID, &req.PositionInColumn)
	return req.ColumnID, req.PositionInColumn, err
}

func (manager *Manager) UpdateTaskDescription(req *taskModels.TaskDescriptionUpdateRequest) (int, int, error) {
	err := manager.Conn.QueryRow("UPDATE tasks SET description = $1 WHERE id = $2 RETURNING column_id, position_in_column", req.Description, req.ID).Scan(&req.ColumnID, &req.PositionInColumn)
	return req.ColumnID, req.PositionInColumn, err
}

func (manager *Manager) UpdateTaskStatus(req *taskModels.TaskStatusUpdateRequest) (int, int, error) {
	err := manager.Conn.QueryRow("UPDATE tasks SET status = $1 WHERE id = $2 RETURNING column_id, position_in_column", req.Status, req.ID).Scan(&req.ColumnID, &req.PositionInColumn)
	return req.ColumnID, req.PositionInColumn, err
}

func (manager *Manager) UpdateTaskDeadline(req *taskModels.TaskDeadlineUpdateRequest) (int, int, error) {
	err := manager.Conn.QueryRow("UPDATE tasks SET deadline = $1 WHERE id = $2 RETURNING column_id, position_in_column", req.Deadline, req.ID).Scan(&req.ColumnID, &req.PositionInColumn)
	return req.ColumnID, req.PositionInColumn, err
}

func (manager *Manager) UpdateTaskDeveloper(req *taskModels.TaskDeveloperUpdateRequest) (int, int, error) {
	err := manager.Conn.QueryRow("UPDATE tasks SET developer_id = $1 WHERE id = $2 RETURNING column_id, position_in_column", req.DeveloperID, req.ID).Scan(&req.ColumnID, &req.PositionInColumn)
	return req.ColumnID, req.PositionInColumn, err
}

func (manager *Manager) UpdateTaskColumn(req *taskModels.TaskColumnUpdateRequest) (int, int, error) {
	err := manager.Conn.QueryRow("UPDATE tasks SET column_id = $1, position_in_column = $2 WHERE id = $3 RETURNING column_id, position_in_column", req.ColumnID, req.PositionInColumn, req.ID).Scan(&req.ColumnID, &req.PositionInColumn)
	return req.ColumnID, req.PositionInColumn, err
}
