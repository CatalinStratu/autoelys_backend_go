ALTER TABLE vehicles
ADD COLUMN user_id BIGINT UNSIGNED NOT NULL AFTER id,
ADD INDEX idx_user_id (user_id),
ADD CONSTRAINT fk_vehicles_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
