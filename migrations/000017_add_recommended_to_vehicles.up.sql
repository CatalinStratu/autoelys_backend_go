ALTER TABLE vehicles ADD COLUMN recommended TINYINT(1) NOT NULL DEFAULT 0 AFTER status;
CREATE INDEX idx_vehicles_recommended ON vehicles(recommended);
