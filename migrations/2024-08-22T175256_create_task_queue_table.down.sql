DROP TRIGGER IF EXISTS set_updated_at ON orchestrator.task_queue;

DROP FUNCTION IF EXISTS update_task_timestamp;

DROP TABLE IF EXISTS orchestrator.task_queue;

DROP TYPE IF EXISTS task_state;
