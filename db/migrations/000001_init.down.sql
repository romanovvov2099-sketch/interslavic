-- Удаление триггера и функции
DROP TRIGGER IF EXISTS update_lesson_progress_updated_at ON lesson_progress;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаление таблиц в обратном порядке (из-за зависимостей)
DROP TABLE IF EXISTS user_answers;
DROP TABLE IF EXISTS lesson_progress;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS lessons;
DROP TABLE IF EXISTS courses;
DROP TABLE IF EXISTS users;