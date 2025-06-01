CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    google_id VARCHAR(255) DEFAULT '',
    first_name VARCHAR(255) NOT NULL DEFAULT '',
    last_name VARCHAR(255) NOT NULL DEFAULT '',
    image_url TEXT DEFAULT '',
    country_code VARCHAR(10) NOT NULL DEFAULT '',
    phone_number VARCHAR(20) NOT NULL UNIQUE DEFAULT '',
    email VARCHAR(255) NOT NULL UNIQUE DEFAULT '',
    password VARCHAR(255) NOT NULL DEFAULT '',
    role int NOT NULL DEFAULT 1,
    status int NOT NULL DEFAULT 1,
    device_token TEXT DEFAULT '',
    email_notify BOOLEAN NOT NULL DEFAULT false,
    push_notify BOOLEAN NOT NULL DEFAULT false,
    created_at bigint NOT NULL DEFAULT EXTRACT (EPOCH FROM CURRENT_TIMESTAMP),
    updated_at bigint NOT NULL DEFAULT EXTRACT (EPOCH FROM CURRENT_TIMESTAMP),
    last_login bigint NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS brands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image_path TEXT DEFAULT '',
    name VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS series (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brand_id UUID NOT NULL REFERENCES brands(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brand_id UUID NOT NULL REFERENCES brands(id) ON DELETE CASCADE,
    series_id UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id),
    receiver_id UUID NOT NULL REFERENCES users(id),
    message TEXT NOT NULL DEFAULT '',
    read BOOLEAN NOT NULL DEFAULT false,
    created_at bigint NOT NULL DEFAULT EXTRACT (EPOCH FROM CURRENT_TIMESTAMP)
);

CREATE TABLE IF NOT EXISTS listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(255) NOT NULL DEFAULT '',
    brand_id UUID NOT NULL REFERENCES brands(id) ON DELETE CASCADE,
    series_id UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,
    model_id UUID NOT NULL REFERENCES models(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    status INT NOT NULL DEFAULT 1,
    currency VARCHAR(3) NOT NULL DEFAULT '',
    price INT NOT NULL DEFAULT 0,
    city VARCHAR(255) NOT NULL DEFAULT '',
    district VARCHAR(255) NOT NULL DEFAULT '',
    neighborhood VARCHAR(255) NOT NULL DEFAULT '',
    images UUID[] DEFAULT '{}',
    fuel_type VARCHAR(50) NOT NULL DEFAULT '',
    transmission_type VARCHAR(50) NOT NULL DEFAULT '',
    body_type VARCHAR(50) NOT NULL DEFAULT '',
    drive_type VARCHAR(50) NOT NULL DEFAULT '',
    engine_power INT NOT NULL DEFAULT 0,
    engine_volume INT NOT NULL DEFAULT 0,
    kilometers INT NOT NULL DEFAULT 0,
    year INT NOT NULL DEFAULT 0,
    color VARCHAR(50) NOT NULL DEFAULT '',
    heavy_damage BOOLEAN NOT NULL DEFAULT false,
    front_bumper VARCHAR(50) NOT NULL DEFAULT '',
    front_hood VARCHAR(50) NOT NULL DEFAULT '',
    roof VARCHAR(50) NOT NULL DEFAULT '',
    front_right_door VARCHAR(50) NOT NULL DEFAULT '',
    rear_right_door VARCHAR(50) NOT NULL DEFAULT '',
    front_left_mudguard VARCHAR(50) NOT NULL DEFAULT '',
    front_left_door VARCHAR(50) NOT NULL DEFAULT '',
    rear_left_door VARCHAR(50) NOT NULL DEFAULT '',
    rear_left_mudguard VARCHAR(50) NOT NULL DEFAULT '',
    rear_bumper VARCHAR(50) NOT NULL DEFAULT '' ,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at bigint NOT NULL DEFAULT EXTRACT (EPOCH FROM CURRENT_TIMESTAMP),
    updated_at bigint NOT NULL DEFAULT EXTRACT (EPOCH FROM CURRENT_TIMESTAMP)
);

CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    path TEXT NOT NULL DEFAULT '',
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at bigint NOT NULL DEFAULT EXTRACT (EPOCH FROM CURRENT_TIMESTAMP)
);

CREATE TABLE IF NOT EXISTS image_predictions (
    image_id UUID PRIMARY KEY REFERENCES images(id) ON DELETE CASCADE,
    prediction BOOLEAN NOT NULL DEFAULT false,
    confidence FLOAT NOT NULL DEFAULT 0.0,
    created_at bigint NOT NULL DEFAULT EXTRACT (EPOCH FROM CURRENT_TIMESTAMP)
);

CREATE TABLE IF NOT EXISTS favorites (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    created_at bigint NOT NULL DEFAULT EXTRACT (EPOCH FROM CURRENT_TIMESTAMP),
    PRIMARY KEY (user_id, listing_id)
);

CREATE TABLE IF NOT EXISTS predicts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at bigint NOT NULL DEFAULT EXTRACT (EPOCH FROM CURRENT_TIMESTAMP),
    brand VARCHAR(255) NOT NULL DEFAULT '',
    series VARCHAR(255) NOT NULL DEFAULT '',
    model VARCHAR(255) NOT NULL DEFAULT '',
    year INT NOT NULL DEFAULT 0,
    mileage FLOAT NOT NULL DEFAULT 0.0,
    engine_volume FLOAT NOT NULL DEFAULT 0.0,
    engine_power FLOAT NOT NULL DEFAULT 0.0,
    accident_history FLOAT NOT NULL DEFAULT 0.0,
    transmission_type VARCHAR(255) NOT NULL DEFAULT '',
    fuel_type VARCHAR(255) NOT NULL DEFAULT '',
    body_type VARCHAR(255) NOT NULL DEFAULT '',
    color VARCHAR(255) NOT NULL DEFAULT '',
    original_parts INT NOT NULL DEFAULT 0,
    replaced_parts INT NOT NULL DEFAULT 0,
    painted_parts INT NOT NULL DEFAULT 0,
    price FLOAT NOT NULL DEFAULT 0.0,
    r2 FLOAT NOT NULL DEFAULT 0.0,
    mae FLOAT NOT NULL DEFAULT 0.0
);
