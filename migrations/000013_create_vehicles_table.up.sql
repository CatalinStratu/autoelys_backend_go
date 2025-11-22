CREATE TABLE IF NOT EXISTS vehicles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    -- General Information
    title VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    description LONGTEXT,
    price DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'lei',
    negotiable BOOLEAN DEFAULT FALSE,

    -- Vehicle Details (Foreign Keys)
    person_type_id TINYINT UNSIGNED NOT NULL,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    engine_capacity INT,
    power_hp INT,
    fuel_type_id TINYINT UNSIGNED NOT NULL,
    body_type_id TINYINT UNSIGNED NOT NULL,
    kilometers INT,
    color VARCHAR(50),
    year INT NOT NULL,
    number_of_keys INT,
    condition_id TINYINT UNSIGNED NOT NULL,
    transmission_id TINYINT UNSIGNED NOT NULL,
    steering_id TINYINT UNSIGNED NOT NULL,
    registered BOOLEAN DEFAULT FALSE,

    -- Location
    city VARCHAR(100) NOT NULL,

    -- Contact Information
    contact_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- Foreign Keys
    FOREIGN KEY (person_type_id) REFERENCES person_types(id),
    FOREIGN KEY (fuel_type_id) REFERENCES fuel_types(id),
    FOREIGN KEY (body_type_id) REFERENCES body_types(id),
    FOREIGN KEY (condition_id) REFERENCES conditions(id),
    FOREIGN KEY (transmission_id) REFERENCES transmissions(id),
    FOREIGN KEY (steering_id) REFERENCES steerings(id),

    -- Indexes
    INDEX idx_person_type_id (person_type_id),
    INDEX idx_brand (brand),
    INDEX idx_model (model),
    INDEX idx_category (category),
    INDEX idx_fuel_type_id (fuel_type_id),
    INDEX idx_body_type_id (body_type_id),
    INDEX idx_year (year),
    INDEX idx_price (price),
    INDEX idx_city (city),
    INDEX idx_condition_id (condition_id),
    INDEX idx_transmission_id (transmission_id),
    INDEX idx_steering_id (steering_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
