 -- Promotions
CREATE TYPE promotion_status AS ENUM ('active', 'inactive', 'expired');

CREATE TABLE promotions (
    promotion_id SERIAL PRIMARY KEY,
    promotion_name VARCHAR(255) NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    status promotion_status DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) PARTITION BY RANGE (start_date);

CREATE INDEX idx_promotion_start_date ON promotions (start_date);
CREATE INDEX idx_promotion_end_date ON promotions (end_date);
CREATE INDEX idx_promotion_status ON promotions (status);

CREATE OR REPLACE FUNCTION update_promotion_status()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.end_date < CURRENT_TIMESTAMP THEN
        NEW.status = 'inactive';  -- Update the status to 'inactive'
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
