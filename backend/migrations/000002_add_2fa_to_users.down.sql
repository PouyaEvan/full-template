ALTER TABLE users DROP COLUMN IF EXISTS is_two_factor_enabled;
ALTER TABLE users DROP COLUMN IF EXISTS two_factor_secret;
ALTER TABLE users DROP COLUMN IF EXISTS two_factor_backup_codes;
