CREATE TABLE IF NOT EXISTS orchestrator.scheduler_lock (
    id UUID PRIMARY KEY,
    instance_id UUID UNIQUE NOT NULL,
    last_heartbeat TIMESTAMP NOT NULL
);

INSERT INTO orchestrator.scheduler_lock (id, instance_id, last_heartbeat)
VALUES (gen_random_uuid(), gen_random_uuid(), '1970-01-01 00:00:00');

