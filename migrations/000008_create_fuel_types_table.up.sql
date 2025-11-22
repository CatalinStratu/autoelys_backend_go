CREATE TABLE IF NOT EXISTS fuel_types (
    id TINYINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Seed data
INSERT INTO fuel_types (name, display_name) VALUES
('benzina', 'Benzină'),
('motorina', 'Motorină'),
('electric', 'Electric'),
('hibrid', 'Hibrid'),
('gpl', 'GPL'),
('hybrid_benzina', 'Hybrid Benzină'),
('hybrid_motorina', 'Hybrid Motorină');
