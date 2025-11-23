ALTER TABLE vehicles
DROP FOREIGN KEY fk_vehicles_user_id,
DROP INDEX idx_user_id,
DROP COLUMN user_id;
