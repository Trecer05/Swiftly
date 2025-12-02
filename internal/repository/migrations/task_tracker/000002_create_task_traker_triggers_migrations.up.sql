
CREATE OR REPLACE FUNCTION reorder_tasks_on_insert()
RETURNS TRIGGER AS $$
BEGIN
    -- Сдвигаем ВСЕ задачи начиная с новой позиции
    UPDATE tasks
    SET position_in_column = position_in_column + 1
    WHERE column_id = NEW.column_id
      AND position_in_column >= NEW.position_in_column;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION reorder_tasks_on_update()
RETURNS TRIGGER AS $$
BEGIN
    -- CASE 1: Задача остаётся в той же колонке
    IF NEW.column_id = OLD.column_id THEN
        
        -- Перемещение вниз (пример: 1 → 4)
        IF NEW.position_in_column > OLD.position_in_column THEN
            UPDATE tasks
            SET position_in_column = position_in_column - 1
            WHERE column_id = NEW.column_id
              AND position_in_column > OLD.position_in_column
              AND position_in_column <= NEW.position_in_column
              AND id <> OLD.id;

        -- Перемещение вверх (пример: 5 → 2)
        ELSIF NEW.position_in_column < OLD.position_in_column THEN
            UPDATE tasks
            SET position_in_column = position_in_column + 1
            WHERE column_id = NEW.column_id
              AND position_in_column < OLD.position_in_column
              AND position_in_column >= NEW.position_in_column
              AND id <> OLD.id;
        END IF;

        RETURN NEW;
    END IF;


    -----------------------------------------------------------------
    -- CASE 2: Перенос задачи в ДРУГУЮ колонку
    -----------------------------------------------------------------

    -- 1. Сжимаем старую колонку, уменьшая позиции ниже удалённой
    UPDATE tasks
    SET position_in_column = position_in_column - 1
    WHERE column_id = OLD.column_id
      AND position_in_column > OLD.position_in_column;

    -- 2. Сдвигаем новую колонку, освобождая позицию
    UPDATE tasks
    SET position_in_column = position_in_column + 1
    WHERE column_id = NEW.column_id
      AND position_in_column >= NEW.position_in_column;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_reorder_tasks_insert
BEFORE INSERT ON tasks
FOR EACH ROW
EXECUTE FUNCTION reorder_tasks_on_insert();

CREATE TRIGGER trg_reorder_tasks_update
BEFORE UPDATE OF column_id, position_in_column ON tasks
FOR EACH ROW
WHEN (OLD.position_in_column IS DISTINCT FROM NEW.position_in_column 
   OR OLD.column_id IS DISTINCT FROM NEW.column_id)
EXECUTE FUNCTION reorder_tasks_on_update();
