-- +goose up
ALTER TABLE users
ADD COLUMN is_red BOOL DEFAULT false;

-- +goose down
ALTER TABLE users
DROP COLUMN is_red;