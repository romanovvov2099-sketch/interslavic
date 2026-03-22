BEGIN;

-- =====================================================
-- Удаление данных в обратном порядке зависимостей
-- =====================================================

-- 1. Сначала удаляем задачи (tasks), так как они ссылаются на уроки (lessons)
DELETE FROM tasks
WHERE lesson_id IN (
    SELECT id FROM lessons
    WHERE course_id IN (
        SELECT id FROM courses WHERE title IN (
            'Модуль 1. Основы общения',
            'Модуль 2. Базовые слова',
            'Модуль 3. Простые фразы',
            'Модуль 4. Практика общения'
        )
    )
);

-- 2. Удаляем уроки (lessons)
DELETE FROM lessons
WHERE course_id IN (
    SELECT id FROM courses WHERE title IN (
        'Модуль 1. Основы общения',
        'Модуль 2. Базовые слова',
        'Модуль 3. Простые фразы',
        'Модуль 4. Практика общения'
    )
);

-- 3. Удаляем курсы (courses)
DELETE FROM courses
WHERE title IN (
    'Модуль 1. Основы общения',
    'Модуль 2. Базовые слова',
    'Модуль 3. Простые фразы',
    'Модуль 4. Практика общения'
);

-- 4. Удаляем пользователей (users), которые были добавлены в up-миграции
DELETE FROM users
WHERE email IN (
    'admin@slavicstudy.com',
    'ivan@slavicstudy.com',
    'maria@slavicstudy.com'
);

-- =====================================================
-- Сброс последовательностей (если используются SERIAL/BIGSERIAL)
-- =====================================================

-- Сброс последовательности для таблицы users
SELECT setval('users_id_seq', COALESCE((SELECT MAX(id) FROM users), 1), false);

-- Сброс последовательности для таблицы courses
SELECT setval('courses_id_seq', COALESCE((SELECT MAX(id) FROM courses), 1), false);

-- Сброс последовательности для таблицы lessons
SELECT setval('lessons_id_seq', COALESCE((SELECT MAX(id) FROM lessons), 1), false);

-- Сброс последовательности для таблицы tasks
SELECT setval('tasks_id_seq', COALESCE((SELECT MAX(id) FROM tasks), 1), false);

COMMIT;