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
SELECT w.id, w.word, w.next_review, s.level, array_agg(t.translation)::text[] as translations
FROM words w
JOIN translations t ON w.id = t.word_id
JOIN stages s ON w.stage_id = s.id
WHERE w.user_id = $1
GROUP BY w.id, s.level;

-- name: GetWordByID :one
SELECT w.id, w.word, w.next_review, s.level, array_agg(t.translation)::text[] as translations
FROM words w
JOIN translations t ON w.id = t.word_id
JOIN stages s ON w.stage_id = s.id
WHERE w.user_id = $1 AND w.id = $2
GROUP BY w.id, s.level;

-- name: GetUserReviewWords :many
SELECT w.id, w.word, w.next_review, s.level, array_agg(t.translation)::text[] as translations
FROM words w
JOIN translations t ON w.id = t.word_id
JOIN stages s ON w.stage_id = s.id
WHERE w.user_id = $1 AND w.next_review <= now()
GROUP BY w.id, s.level;

-- name: GetUserReviewWordsCount :one
SELECT count(w.id)
FROM words w
WHERE w.user_id = $1 AND w.next_review <= now();

-- name: GetUserReviewsCountInRange :many
SELECT count(id), next_review
FROM words
WHERE user_id = $1 AND next_review > $2 AND next_review < $3
GROUP BY next_review
ORDER BY next_review ASC;

-- name: UpdateWordStage :exec
UPDATE words
SET stage_id = $1, next_review = $2
WHERE id = $3 AND user_id = $4;