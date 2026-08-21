-- init-db.sql

-- 1. Создаем таблицу users
-- Обратите внимание: city имеет дефолтное значение пустой строки, но в INSERT мы его явно укажем
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    age INTEGER NOT NULL,
    experience INTEGER,
    city VARCHAR(100) NOT NULL DEFAULT ''
);

-- 2. Создаем таблицу items
-- user_id ссылается на users(id) с каскадным удалением
-- price имеет проверку: цена не может быть отрицательной
-- status по умолчанию 'active'
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(50) NOT NULL,
    description VARCHAR(200) NOT NULL DEFAULT '',
    price NUMERIC(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    CONSTRAINT items_price_check CHECK (price >= 0),
    CONSTRAINT items_user_id_fkey FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE
);

-- 3. Наполняем таблицу users тестовыми данными
-- Важно: все поля NOT NULL, поэтому нельзя оставлять их пустыми
INSERT INTO users (name, age, experience, city) VALUES 
('Алексей', 28, 5, 'Москва'),
('Мария', 32, 8, 'Санкт-Петербург'),
('Дмитрий', 25, 2, 'Казань')
ON CONFLICT DO NOTHING;

-- 4. Наполняем таблицу items тестовыми данными
-- Важно: user_id должен существовать в таблице users. 
-- Мы используем ID, которые только что вставили (1, 2, 3).
-- Цена должна быть >= 0.
INSERT INTO items (user_id, title, description, price, status) VALUES 
(1, 'Ноутбук Pro', 'Мощный ноутбук для разработки', 120000.00, 'active'),
(1, 'Механическая клавиатура', 'Тихие свитчи, RGB подсветка', 15000.00, 'active'),
(2, 'Графический планшет', 'Для ретуши и рисования', 45000.00, 'active'),
(3, 'Беспроводные наушники', 'Шумоподавление, долгая батарея', 12000.00, 'active')
ON CONFLICT DO NOTHING;