BEGIN;

-- =====================================================
-- USERS
-- =====================================================

INSERT INTO users (fullname, email, login, password, role, last_login)
VALUES
    ('Admin', 'admin@slavicstudy.com', 'admin', '$2a$10$N9qo8uLOickgx2ZMRZoMy.Mr/AJ5R2sF5qR5qR5qR5qR5qR5qR5q', 'admin', CURRENT_TIMESTAMP),
    ('Ivan Petrov', 'ivan@slavicstudy.com', 'ivanpetrov', '$2a$10$examplehash111111111111111111111111111111111111111', 'student', CURRENT_TIMESTAMP - INTERVAL '1 day'),
    ('Maria Sokolova', 'maria@slavicstudy.com', 'mariasokolova', '$2a$10$examplehash222222222222222222222222222222222222222', 'student', CURRENT_TIMESTAMP - INTERVAL '2 days')
ON CONFLICT (email) DO NOTHING;

-- =====================================================
-- COURSES / MODULES
-- =====================================================

INSERT INTO courses (title, description, theory)
SELECT *
FROM (
    VALUES
    (
        'Модуль 1. Основы общения',
        'Первый модуль для освоения приветствий, знакомства и простых ответов.',
        'Практика базового общения на межславянском языке.'
    ),
    (
        'Модуль 2. Базовые слова',
        'Модуль с лексикой по темам: семья, дом, люди, повседневная жизнь.',
        'Практика употребления базовой лексики в коротких фразах.'
    ),
    (
        'Модуль 3. Простые фразы',
        'Модуль по построению простых предложений и описанию себя.',
        'Практика составления коротких предложений на межславянском языке.'
    ),
    (
        'Модуль 4. Практика общения',
        'Модуль с мини-диалогами и прикладными речевыми ситуациями.',
        'Практика общения в реальных бытовых сценариях.'
    )
) AS v(title, description, theory)
WHERE NOT EXISTS (
    SELECT 1 FROM courses c WHERE c.title = v.title
);

-- =====================================================
-- LESSONS
-- =====================================================

-- Модуль 1
INSERT INTO lessons (course_id, title, content, multimedia, position)
SELECT c.id, l.title, l.content, NULL, l.position
FROM courses c
JOIN (
    VALUES
    ('Приветствия', 'Практика базовых приветствий.', 1),
    ('Знакомство', 'Практика фраз для знакомства.', 2),
    ('Кто ты?', 'Практика ответа на вопрос о себе.', 3),
    ('Простые ответы', 'Практика коротких ответов в диалоге.', 4)
) AS l(title, content, position)
ON c.title = 'Модуль 1. Основы общения'
WHERE NOT EXISTS (
    SELECT 1 FROM lessons ls WHERE ls.course_id = c.id AND ls.title = l.title
);

-- Модуль 2
INSERT INTO lessons (course_id, title, content, multimedia, position)
SELECT c.id, l.title, l.content, NULL, l.position
FROM courses c
JOIN (
    VALUES
    ('Семья', 'Практика слов по теме семьи.', 1),
    ('Дом', 'Практика слов по теме дома.', 2),
    ('Люди и профессии', 'Практика слов о людях и профессиях.', 3),
    ('Повседневные слова', 'Практика часто употребляемых слов.', 4)
) AS l(title, content, position)
ON c.title = 'Модуль 2. Базовые слова'
WHERE NOT EXISTS (
    SELECT 1 FROM lessons ls WHERE ls.course_id = c.id AND ls.title = l.title
);

-- Модуль 3
INSERT INTO lessons (course_id, title, content, multimedia, position)
SELECT c.id, l.title, l.content, NULL, l.position
FROM courses c
JOIN (
    VALUES
    ('Местоимения', 'Практика местоимений в речи.', 1),
    ('Глагол byti', 'Практика глагола byti.', 2),
    ('Простое предложение', 'Практика построения предложения.', 3),
    ('Описание себя', 'Практика короткого рассказа о себе.', 4)
) AS l(title, content, position)
ON c.title = 'Модуль 3. Простые фразы'
WHERE NOT EXISTS (
    SELECT 1 FROM lessons ls WHERE ls.course_id = c.id AND ls.title = l.title
);

-- Модуль 4
INSERT INTO lessons (course_id, title, content, multimedia, position)
SELECT c.id, l.title, l.content, NULL, l.position
FROM courses c
JOIN (
    VALUES
    ('Приветствие в диалоге', 'Практика приветствия в коротком диалоге.', 1),
    ('Представление себя', 'Практика самопрезентации.', 2),
    ('Вопрос и ответ', 'Практика ответов на простые вопросы.', 3),
    ('Мини-диалог', 'Практика завершённого короткого диалога.', 4)
) AS l(title, content, position)
ON c.title = 'Модуль 4. Практика общения'
WHERE NOT EXISTS (
    SELECT 1 FROM lessons ls WHERE ls.course_id = c.id AND ls.title = l.title
);

-- =====================================================
-- TASKS
-- =====================================================

-- -----------------------------
-- МОДУЛЬ 1. УРОК 1. Приветствия
-- самый простой уровень
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери перевод: Dobry den', 'Добрый день', ARRAY['До свидания', 'Добрый день', 'Спасибо', 'Спокойной ночи']::TEXT[]),
    (1, 'Выбери приветствие: Доброе утро', 'Dobro utro', ARRAY['Dobro utro', 'Dobry večer', 'Hvala', 'Do viděnja']::TEXT[]),
    (2, 'Напиши на межславянском: Добрый день', 'Dobry den', ARRAY[]::TEXT[]),
    (2, 'Напиши на межславянском: Доброе утро', 'Dobro utro', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Приветствия'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

-- -----------------------------
-- МОДУЛЬ 1. УРОК 2. Знакомство
-- сложнее: короткие фразы
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери перевод: Ja jesm Ivan', 'Я Иван', ARRAY['Я Иван', 'Ты Иван', 'Он Иван', 'Мы Иван']::TEXT[]),
    (1, 'Выбери правильную фразу для знакомства', 'Ja jesm Maria', ARRAY['Dobry den', 'Ja jesm Maria', 'Do viděnja', 'Hvala']::TEXT[]),
    (2, 'Напиши на межславянском: Я Мария', 'Ja jesm Maria', ARRAY[]::TEXT[]),
    (2, 'Напиши на межславянском: Я студент', 'Ja jesm student', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Знакомство'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

-- -----------------------------
-- МОДУЛЬ 1. УРОК 3. Кто ты?
-- сложнее: вопрос-ответ
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери ответ на вопрос: Kto ty?', 'Ja jesm student', ARRAY['Ja jesm student', 'Dobry den', 'Moj dom', 'Ty jesi učitelj']::TEXT[]),
    (1, 'Выбери перевод: Kto ty?', 'Кто ты?', ARRAY['Как тебя зовут?', 'Где ты?', 'Кто ты?', 'Что это?']::TEXT[]),
    (2, 'Ответь на межславянском: Кто ты?', 'Ja jesm student', ARRAY[]::TEXT[]),
    (2, 'Напиши на межславянском: Я учитель', 'Ja jesm učitelj', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Кто ты?'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

-- -----------------------------
-- МОДУЛЬ 1. УРОК 4. Простые ответы
-- ещё сложнее: мини-реплики
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери подходящий ответ на "Dobry den"', 'Dobry den', ARRAY['Dobry den', 'Ja dom', 'Student', 'Kto ty']::TEXT[]),
    (1, 'Выбери подходящий ответ на "Kto ty?"', 'Ja jesm Ivan', ARRAY['Dobry den', 'Ja jesm Ivan', 'Do viděnja', 'Dobro utro']::TEXT[]),
    (2, 'Ответь на реплику: Dobry den', 'Dobry den', ARRAY[]::TEXT[]),
    (2, 'Ответь на вопрос: Kto ty?', 'Ja jesm student', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Простые ответы'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

-- -----------------------------
-- МОДУЛЬ 2. УРОК 1. Семья
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери перевод: matka', 'мать', ARRAY['отец', 'мать', 'сын', 'дом']::TEXT[]),
    (1, 'Выбери слово по теме "семья"', 'brat', ARRAY['brat', 'stol', 'okno', 'pisati']::TEXT[]),
    (2, 'Напиши на межславянском: брат', 'brat', ARRAY[]::TEXT[]),
    (2, 'Напиши на межславянском: мать', 'matka', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Семья'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

-- -----------------------------
-- МОДУЛЬ 2. УРОК 2. Дом
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери перевод: dom', 'дом', ARRAY['семья', 'окно', 'дом', 'город']::TEXT[]),
    (1, 'Выбери слово по теме "дом"', 'okno', ARRAY['okno', 'učitelj', 'matka', 'ja']::TEXT[]),
    (2, 'Напиши на межславянском: окно', 'okno', ARRAY[]::TEXT[]),
    (2, 'Напиши на межславянском: Это дом', 'To jest dom', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Дом'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

-- -----------------------------
-- МОДУЛЬ 2. УРОК 3. Люди и профессии
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери перевод: učitelj', 'учитель', ARRAY['врач', 'учитель', 'студент', 'брат']::TEXT[]),
    (1, 'Выбери правильную фразу', 'On jest učitelj', ARRAY['On jest učitelj', 'Dom jest brat', 'Ja jest okno', 'Ty jest matka']::TEXT[]),
    (2, 'Напиши на межславянском: Он учитель', 'On jest učitelj', ARRAY[]::TEXT[]),
    (2, 'Напиши на межславянском: Я студент', 'Ja jesm student', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Люди и профессии'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

-- -----------------------------
-- МОДУЛЬ 2. УРОК 4. Повседневные слова
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери перевод: kniga', 'книга', ARRAY['ручка', 'дом', 'книга', 'стол']::TEXT[]),
    (1, 'Выбери подходящую фразу', 'To jest kniga', ARRAY['To jest kniga', 'Ja jesm kniga', 'Ty dom', 'On brat jest dom']::TEXT[]),
    (2, 'Напиши на межславянском: Это книга', 'To jest kniga', ARRAY[]::TEXT[]),
    (2, 'Напиши на межславянском: Это стол', 'To jest stol', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Повседневные слова'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

-- -----------------------------
-- МОДУЛЬ 3. УРОК 1. Местоимения
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери перевод: my', 'мы', ARRAY['я', 'вы', 'мы', 'они']::TEXT[]),
    (1, 'Выбери местоимение для "они"', 'oni', ARRAY['on', 'ona', 'oni', 'my']::TEXT[]),
    (2, 'Напиши на межславянском: мы', 'my', ARRAY[]::TEXT[]),
    (2, 'Напиши на межславянском: они', 'oni', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Местоимения'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

-- -----------------------------
-- МОДУЛЬ 3. УРОК 2. Глагол byti
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери правильную фразу: Я студент', 'Ja jesm student', ARRAY['Ja jesm student', 'Ja student dom', 'Student ja to', 'Ja byti student']::TEXT[]),
    (1, 'Выбери правильную фразу: Он учитель', 'On jest učitelj', ARRAY['On jest učitelj', 'On jesm učitelj', 'On dom učitelj', 'On ty učitelj']::TEXT[]),
    (2, 'Напиши на межславянском: Мы студенты', 'My jesmo studenti', ARRAY[]::TEXT[]),
    (2, 'Напиши на межславянском: Она учитель', 'Ona jest učitelj', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Глагол byti'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

-- -----------------------------
-- МОДУЛЬ 3. УРОК 3. Простое предложение
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери правильный перевод: Это мой дом', 'To jest moj dom', ARRAY['To jest moj dom', 'Ja jesm dom moj', 'Moj to jest ja', 'Dom jest ty']::TEXT[]),
    (1, 'Выбери правильный перевод: Это моя книга', 'To jest moja kniga', ARRAY['To jest moja kniga', 'Ja moja kniga', 'Kniga jest ona', 'Moja ja kniga']::TEXT[]),
    (2, 'Напиши на межславянском: Это мой брат', 'To jest moj brat', ARRAY[]::TEXT[]),
    (2, 'Напиши на межславянском: Это моя мать', 'To jest moja matka', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Простое предложение'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

-- -----------------------------
-- МОДУЛЬ 3. УРОК 4. Описание себя
-- самый сложный внутри модуля
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери правильную самопрезентацию', 'Ja jesm Ivan. Ja jesm student.', ARRAY[
        'Ja jesm Ivan. Ja jesm student.',
        'Ivan student dom.',
        'Ja dom student.',
        'Student jest ja ty.'
    ]::TEXT[]),
    (1, 'Выбери лучший ответ на вопрос "Kto ty?"', 'Ja jesm Maria. Ja jesm učitelj.', ARRAY[
        'Ja jesm Maria. Ja jesm učitelj.',
        'Dobry den okno.',
        'To jest brat.',
        'My dom.'
    ]::TEXT[]),
    (2, 'Напиши 2 короткие фразы о себе на межславянском.', 'Ja jesm Ivan. Ja jesm student.', ARRAY[]::TEXT[]),
    (2, 'Напиши: Меня зовут Мария. Я учитель.', 'Ja jesm Maria. Ja jesm učitelj.', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Описание себя'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

-- -----------------------------
-- МОДУЛЬ 4. УРОК 1. Приветствие в диалоге
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери лучший ответ в диалоге: — Dobry den! — ?', 'Dobry den!', ARRAY['Dobry den!', 'Ja dom.', 'Student.', 'To kniga.']::TEXT[]),
    (2, 'Ответь на реплику: Dobro utro!', 'Dobro utro!', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Приветствие в диалоге'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

-- -----------------------------
-- МОДУЛЬ 4. УРОК 2. Представление себя
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери правильный ответ: — Kto ty? — ?', 'Ja jesm student.', ARRAY['Ja jesm student.', 'Dobry večer.', 'To jest dom.', 'Brat okno.']::TEXT[]),
    (2, 'Ответь на вопрос: Kto ty?', 'Ja jesm učitelj.', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Представление себя'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

-- -----------------------------
-- МОДУЛЬ 4. УРОК 3. Вопрос и ответ
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери ответ: — Kto to jest? — ?', 'To jest moj brat.', ARRAY['To jest moj brat.', 'Ja jesm brat.', 'Dobry den.', 'My studenti.']::TEXT[]),
    (2, 'Ответь: — Kto to jest? — Это моя мать.', 'To jest moja matka.', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Вопрос и ответ'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

-- -----------------------------
-- МОДУЛЬ 4. УРОК 4. Мини-диалог
-- самый сложный блок
-- -----------------------------
INSERT INTO tasks (lesson_id, task_type, question, answer, choises)
SELECT l.id, t.task_type, t.question, t.answer, t.choises
FROM lessons l
JOIN (
    VALUES
    (1, 'Выбери завершённый диалог: — Dobry den! — ? — Kto ty? — ?', 'Dobry den! Ja jesm Ivan.', ARRAY[
        'Dobry den! Ja jesm Ivan.',
        'Dom kniga student.',
        'To jest matka dom.',
        'Ja ty oni my.'
    ]::TEXT[]),
    (2, 'Напиши мини-диалог из 2 реплик: приветствие и представление.', 'Dobry den! Ja jesm Maria.', ARRAY[]::TEXT[]),
    (2, 'Напиши ответ на диалог: — Dobry den! Kto ty?', 'Dobry den! Ja jesm student.', ARRAY[]::TEXT[])
) AS t(task_type, question, answer, choises)
ON l.title = 'Мини-диалог'
WHERE NOT EXISTS (
    SELECT 1 FROM tasks tt WHERE tt.lesson_id = l.id AND tt.question = t.question
);

COMMIT;