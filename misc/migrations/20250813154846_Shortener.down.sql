-- Migration Down: Reverse the URL shortener table creation

-- Drop the trigger first (depends on function)
DROP TRIGGER update_urls_updated_at ON urls;

-- Drop the function
DROP FUNCTION update_url_entity();

-- Drop indexes (they will also be dropped with table, but explicit is better)
DROP INDEX idx_urls_short_code;

-- Drop the main table (this will also drop all associated indexes and constraints)
DROP TABLE urls;
