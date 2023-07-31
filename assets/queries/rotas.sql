-- name: FindRotaByID :one
SELECT ROTAS.*
FROM ROTAS
WHERE ID = $1;

-- name: ListRotasByChannel :many
SELECT ROTAS.*
from ROTAS
WHERE ROTAS.CHANNEL_ID = $1
  AND ROTAS.TEAM_ID = $2;

-- name: SaveRota :one
INSERT INTO ROTAS (TEAM_ID, CHANNEL_ID, NAME, METADATA)
VALUES ($1, $2, $3, $4) RETURNING ID;

-- name: UpdateRota :one
UPDATE ROTAS
SET NAME       = $1,
    METADATA   = $2
WHERE ID = $3
RETURNING ID;