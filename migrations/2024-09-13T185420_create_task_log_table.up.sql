CREATE TABLE IF NOT EXISTS orchestrator.task_logs
(
    id      UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    task_id UUID  NOT NULL,
    log     JSONB NOT NULL,
    CONSTRAINT fk_task
        FOREIGN KEY (task_id)
            REFERENCES orchestrator.task_queue (id)
            ON DELETE CASCADE
);