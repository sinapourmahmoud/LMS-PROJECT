CREATE TABLE IF NOT EXISTS courses(
    id UUID PRIMARY KEY,
    description text NOT NULL,
    teacher_id UUID NOT NULL REFERENCES users(id),
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','published')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ

);
