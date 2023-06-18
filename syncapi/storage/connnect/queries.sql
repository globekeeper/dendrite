-- name: InsertMultiRoomData :one
INSERT INTO syncapi_multiroom_data (
    user_id,
    type,
    data
) VALUES (
    $1,
    $2,
    $3
) ON CONFLICT (user_id, type) DO UPDATE SET id = nextval('syncapi_multiroom_id'), data = $3, ts = current_timestamp
RETURNING id;

-- name: InsertMultiRoomVisibility :exec
INSERT INTO syncapi_multiroom_visibility (
    user_id,
    type,
    room_id,
    expire_ts
) VALUES (
    $1,
    $2,
    $3,
    $4
) ON CONFLICT (user_id, type, room_id) DO UPDATE SET expire_ts = $4;

-- name: SelectMultiRoomVisibilityRooms :many
SELECT room_id FROM syncapi_multiroom_visibility
WHERE user_id = $1 
AND expire_ts > $2;

-- name: SelectMaxId :one
SELECT MAX(id) FROM syncapi_multiroom_data;

-- name: DeleteMultiRoomVisibility :exec
DELETE FROM syncapi_multiroom_visibility
WHERE user_id = $1
AND type = $2
AND room_id = $3;

-- name: DeleteMultiRoomVisibilityByExpireTS :execrows
DELETE FROM syncapi_multiroom_visibility
WHERE expire_ts <= $1;

-- name: UpsertDataRetention :exec
INSERT INTO data_retention (
    space_id,
    timeframe,
    at
) VALUES (
    $1,
    $2,
    $3
) ON CONFLICT (space_id) DO UPDATE SET timeframe = $2, at = $3;

-- name: DeleteDataRetention :exec
DELETE FROM data_retention
WHERE space_id = $1;

-- name: SelectDataRetentions :many
SELECT * FROM data_retention;