-- name: CreateUrlMapping :execresult
INSERT INTO urls (
    original_url,
    short_code,
    expires_at,
    creator_ip
) VALUES (
    $1, $2, $3, $4
);

-- name: GetShortCodeFromOriginalUrl :one
SELECT short_code FROM urls WHERE original_url = $1 LIMIT 1;

-- name: GetOriginalUrlFromShortCode :one
SELECT original_url FROM urls WHERE short_code = $1 LIMIT 1;