-- Remove status column from vehicles table

ALTER TABLE vehicles
DROP INDEX idx_status,
DROP COLUMN status;
