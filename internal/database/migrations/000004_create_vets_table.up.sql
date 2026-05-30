CREATE TABLE IF NOT EXISTS vets (
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID         NOT NULL UNIQUE REFERENCES users(id),
    description     TEXT,
    clinic_name     VARCHAR(255),
    consultation_fee DECIMAL(10,2),
    max_concurrent  INT          NOT NULL DEFAULT 1,
    status          VARCHAR(20)  NOT NULL DEFAULT 'offline',
    location        GEOGRAPHY(Point, 4326),
    rating_avg      DECIMAL(2,1) NOT NULL DEFAULT 0.0,
    reviews_count   INT          NOT NULL DEFAULT 0,
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_vets_status ON vets(status);
CREATE INDEX IF NOT EXISTS idx_vets_location ON vets USING GIST (location);
