CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    image_url TEXT,
    country_code VARCHAR(10) NOT NULL,
    phone_number VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'Regular',
    status VARCHAR(50) NOT NULL DEFAULT 'Active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP NULL,
    UNIQUE (id, phone_number, email)
);

CREATE TABLE IF NOT EXISTS brands (
    id SERIAL PRIMARY KEY,
    logo TEXT,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS series (
    id SERIAL PRIMARY KEY,
    brand_id INT,
    name VARCHAR(255) NOT NULL,
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

CREATE TABLE IF NOT EXISTS models (
    id SERIAL PRIMARY KEY,
    series_id INT,
    name VARCHAR(255) NOT NULL,
    FOREIGN KEY (series_id) REFERENCES series(id)
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'fuel_type') THEN
        CREATE TYPE fuel_type AS ENUM (
            'Diesel',
            'Petrol',
            'Petrol & LPG',
            'Hybrid',
            'Electric'
        );
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'transmission') THEN
        CREATE TYPE transmission AS ENUM (
            'Automatic',
            'Manual',
            'Semiautomatic'
        );
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'body_type') THEN
        CREATE TYPE body_type AS ENUM (
            'Sedan',
            'Hatchback/3',
            'Hatchback/5',
            'Coupe',
            'Cabrio',
            'MPV',
            'Pick-up',
            'Roadster',
            'Station wagon',
            'SUV'
        );
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'part_condition') THEN
        CREATE TYPE part_condition AS ENUM (
            'Original',
            'Painted',
            'Changed'
        );
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'seller_type') THEN
        CREATE TYPE seller_type AS ENUM (
            'Individual',
            'Dealer'
        );
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'currency') THEN
        CREATE TYPE currency AS ENUM (
            'TRY',
            'USD',
            'EUR'
        );
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'drive_type') THEN
        CREATE TYPE drive_type AS ENUM (
            'Front-Wheel Drive',
            'Rear-Wheel Drive',
            'Four-Wheel Drive',
            'All-Wheel Drive'
        );
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS cars (
    id VARCHAR(255) PRIMARY KEY,
    owner_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    currency currency NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    city VARCHAR(100) NOT NULL,
    district VARCHAR(100) NOT NULL,
    neighborhood VARCHAR(100),
    listing_number VARCHAR(50) UNIQUE NOT NULL,
    listing_date TIMESTAMP NOT NULL,
    brand_id INT NOT NULL,
    series_id INT NOT NULL,
    model_id INT NOT NULL,
    year INT NOT NULL CHECK (year >= 1886),
    fuel_type fuel_type NOT NULL,
    transmission transmission NOT NULL,
    mileage INT NOT NULL CHECK (mileage >= 0),
    body_type body_type NOT NULL,
    engine_power INT CHECK (engine_power >= 0),
    engine_volume INT CHECK (engine_volume >= 0),
    drive_type drive_type,
    color VARCHAR(50),
    warranty BOOLEAN DEFAULT FALSE,
    heavy_damage BOOLEAN DEFAULT FALSE,
    seller_type seller_type NOT NULL,
    trade_option BOOLEAN DEFAULT FALSE,
    front_bumper part_condition,
    front_hood part_condition,
    roof part_condition,
    front_right_door part_condition,
    rear_right_door part_condition,
    front_left_mudguard part_condition,
    front_left_door part_condition,
    rear_left_door part_condition,
    rear_left_mudguard part_condition,
    rear_bumper part_condition
);
