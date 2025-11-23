-- Add status column to vehicles table
-- Status values: 1 = active, 2 = inactive, 3 = banned

ALTER TABLE vehicles
ADD COLUMN status TINYINT UNSIGNED NOT NULL DEFAULT 1
COMMENT '1=active, 2=inactive, 3=banned' AFTER user_id,
ADD INDEX idx_status (status);

-- Update existing vehicles to active status
UPDATE vehicles SET status = 1 WHERE status IS NULL OR status = 0;
