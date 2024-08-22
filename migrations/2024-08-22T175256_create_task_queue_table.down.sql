DROP TRIGGER IF EXISTS set_updated_at ON orchestrator.tasks;

DROP FUNCTION IF EXISTS update_task_timestamp;

DROP TABLE IF EXISTS orchestrator.tasks;

DROP TYPE IF EXISTS task_state;
