-- name: CreateUrl :one
INSERT INTO urls (url, short_code)
VALUES ($1, $2)
RETURNING id, url, short_code, created_at, updated_at;

-- name: GetUrl :one
SELECT id, url, short_code, created_at, updated_at
FROM urls
WHERE short_code = $1
LIMIT 1;

-- name: DeleteUrl :one
DELETE
FROM urls
WHERE short_code = $1
RETURNING id;

-- name: StatUrls :one
SELECT id, url, short_code, created_at, updated_at, access_count
FROM urls
WHERE short_code = $1
LIMIT 1;

-- name: UpdateUrl :one
UPDATE urls
SET url = $1, access_count = 0, updated_at = CURRENT_TIMESTAMP
WHERE short_code = $2
RETURNING id, url, short_code, created_at, updated_at;

-- name: IncrementAccessCount :one
UPDATE urls
SET access_count = access_count + 1
WHERE short_code = $1
RETURNING access_count;
