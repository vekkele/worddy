-- name: AddWord :one
INSERT INTO words (word, next_review, stage_id, user_id)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: AddTranslation :exec
INSERT INTO translations (translation, word_id)
VALUES ($1, $2);

-- name: GetStageByLevel :one
SELECT id, level, hours_to_next FROM stages WHERE level = $1;

-- name: GetUserWords :many
SELECT w.id, w.word, w.next_review, s.level, string_agg(t.translation, ', ') as translations
FROM words w
JOIN translations t ON w.id = t.word_id
JOIN stages s ON w.stage_id = s.id
WHERE w.user_id = $1
GROUP BY w.id, s.level;