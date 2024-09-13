CREATE TABLE IF NOT EXISTS orchestrator.task_logs
(
    id  UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    log JSONB NOT NULL
);