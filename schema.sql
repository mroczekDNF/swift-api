DROP TABLE IF EXISTS swift_codes;

CREATE TABLE swift_codes (
    id SERIAL PRIMARY KEY,
    swift_code VARCHAR(11) UNIQUE NOT NULL,
    bank_name TEXT NOT NULL,
    address TEXT,
    country_iso2 CHAR(2) NOT NULL,
    country_name TEXT NOT NULL,
    is_headquarter BOOLEAN NOT NULL,
    headquarter_id INT
);

-- Dodanie indeksu na headquarter_id dla szybkiego wyszukiwania branchy
CREATE INDEX IF NOT EXISTS idx_headquarter_id ON swift_codes (headquarter_id);
