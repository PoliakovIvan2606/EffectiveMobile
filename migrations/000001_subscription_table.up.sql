CREATE TABLE IF NOT EXISTS subscription (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL,
    price INT NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);