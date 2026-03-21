-- Создание таблицы users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    login VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'student',
    reg_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP WITH TIME ZONE
);

-- Создание индексов для users
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_login ON users(login);
CREATE INDEX idx_users_role ON users(role);

-- Создание таблицы courses
CREATE TABLE IF NOT EXISTS courses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    theory TEXT
);

-- Создание таблицы lessons
CREATE TABLE IF NOT EXISTS lessons (
    id SERIAL PRIMARY KEY,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    multimedia TEXT,
    position INTEGER DEFAULT 0
);

-- Создание индексов для lessons
CREATE INDEX idx_lessons_course_id ON lessons(course_id);
CREATE INDEX idx_lessons_position ON lessons(position);

-- Создание таблицы tasks
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    lesson_id INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    task_type INTEGER NOT NULL CHECK (task_type IN (1, 2)), -- 1: Choices, 2: Write
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    choises TEXT[] DEFAULT '{}'
);

-- Создание индексов для tasks
CREATE INDEX idx_tasks_lesson_id ON tasks(lesson_id);
CREATE INDEX idx_tasks_task_type ON tasks(task_type);

-- Создание таблицы lesson_progress
CREATE TABLE IF NOT EXISTS lesson_progress (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'not_started',
    score INTEGER DEFAULT 0,
    completion_date TIMESTAMP WITH TIME ZONE,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, lesson_id)
);

-- Создание индексов для lesson_progress
CREATE INDEX idx_lesson_progress_user_id ON lesson_progress(user_id);
CREATE INDEX idx_lesson_progress_lesson_id ON lesson_progress(lesson_id);
CREATE INDEX idx_lesson_progress_status ON lesson_progress(status);
CREATE INDEX idx_lesson_progress_completion_date ON lesson_progress(completion_date);

-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггер для обновления updated_at в lesson_progress
CREATE TRIGGER update_lesson_progress_updated_at
    BEFORE UPDATE ON lesson_progress
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Вставка начальных данных (опционально)
INSERT INTO users (fullname, email, login, password, role) VALUES
    ('Admin', 'admin@slavicstudy.com', 'admin', '$2a$10$N9qo8uLOickgx2ZMRZoMy.Mr/AJ5R2sF5qR5qR5qR5qR5qR5qR5q', 'admin')
ON CONFLICT (email) DO NOTHING;