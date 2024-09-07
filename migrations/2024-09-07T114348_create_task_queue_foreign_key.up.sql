ALTER TABLE orchestrator.task_queue
    ADD CONSTRAINT fk_task_name
        FOREIGN KEY (name) REFERENCES orchestrator.tasks (name)
            ON DELETE CASCADE;
