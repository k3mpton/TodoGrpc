-- +goose Up
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY, 
    title VARCHAR(255) NOT NULL,
    description TEXT, 
    status BOOLEAN DEFAULT FALSE, 
    due_date TIMESTAMP, 
    user_id bigint not null, 
    created_at TIMESTAMP DEFAULT Now(), 
    updated_at TIMESTAMP DEFAULT Now() 
);

CREATE INDEX idx_tasks_title ON tasks(title);

-- +goose Down
DROP INDEX IF EXISTS idx_tasks_title;
DROP TABLE IF EXISTS tasks; 