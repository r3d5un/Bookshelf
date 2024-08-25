CREATE TYPE task_state AS ENUM ('waiting', 'running', 'complete', 'stopped', 'error');

CREATE TABLE IF NOT EXISTS orchestrator.tasks
(
    id         UUID                 DEFAULT gen_random_uuid() PRIMARY KEY,
    name       VARCHAR(32) NOT NULL,
    state      task_state  NOT NULL DEFAULT 'waiting',
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    run_at     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    task_data  JSONB       NULL     DEFAULT NULL
);

-- Function and trigger to always update the updated_at column when changed
CREATE OR REPLACE FUNCTION update_task_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at
    BEFORE UPDATE
    ON orchestrator.tasks
    FOR EACH ROW
EXECUTE FUNCTION update_task_timestamp();
