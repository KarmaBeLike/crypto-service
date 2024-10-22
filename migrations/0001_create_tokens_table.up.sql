-- Сначала создаем таблицу tokens
CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    price_usd NUMERIC(18, 8) NOT NULL
);

-- Затем создаем таблицу token_price_history
CREATE TABLE token_price_history (
    id SERIAL PRIMARY KEY,
    token_id INTEGER REFERENCES tokens(id),  
    price DECIMAL(18, 8), 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  
);

ALTER TABLE token_price_history ADD COLUMN symbol VARCHAR(10);
