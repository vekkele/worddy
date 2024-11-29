-- name: AddReviewWord :exec
INSERT INTO review_words (word_id, user_id) 
  VALUES ($1, $2)
  ON CONFLICT DO NOTHING;

-- name: CommitWrongReviewAnswer :exec
UPDATE review_words
  SET wrong_answers = wrong_answers + 1
  WHERE word_id = $1 AND user_id = $2;

-- name: DeleteReviewWord :exec
DELETE FROM review_words
  WHERE word_id = $1 AND user_id = $2;

-- name: GetReviewWrongAnswers :one
SELECT wrong_answers FROM review_words
  WHERE word_id = $1 AND user_id = $2;

-- name: GetNextReviewWord :one
SELECT DISTINCT ON (wrong_answers) wrong_answers, word_id, user_id
  FROM review_words
  WHERE user_id = $1
  ORDER BY wrong_answers
  LIMIT 1;
