CREATE TABLE IF NOT EXISTS body_types (
    id TINYINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Seed data
INSERT INTO body_types (name, display_name) VALUES
('sedan', 'Sedan'),
('suv', 'SUV'),
('break', 'Break'),
('coupe', 'Coupe'),
('cabrio', 'Cabrio'),
('hatchback', 'Hatchback'),
('pickup', 'Pickup'),
('van', 'Van'),
('monovolum', 'Monovolum');
