CREATE TABLE IF NOT EXISTS automobiles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    url_hash VARCHAR(191) NOT NULL,
    url LONGTEXT NOT NULL,
    brand_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(191) NOT NULL,
    description LONGTEXT,
    press_release LONGTEXT,
    photos LONGTEXT,
    created_at TIMESTAMP NULL DEFAULT NULL,
    updated_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_url_hash (url_hash),
    INDEX idx_brand_id (brand_id),
    INDEX idx_name (name),
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
