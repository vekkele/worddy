-- name: AddWord :exec
INSERT INTO words (word, next_review, stage_id, user_id)
VALUES ($1, $2, $3, $4);

-- name: AddTranslation :exec
INSERT INTO translations (translation, word_id)
VALUES ($1, $2);

-- name: GetStageByLevel :one
SELECT id, level, hours_to_next FROM stages WHERE level = $1;
