CREATE TABLE IF NOT EXISTS books.series (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    description TEXT NULL,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL
);
