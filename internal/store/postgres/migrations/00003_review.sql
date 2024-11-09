-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS review_words (
  wrong_answers integer NOT NULL DEFAULT 0,
  word_id bigint NOT NULL,
  constraint fk_word
    foreign key (word_id)
    references words(id)
    on delete cascade,
  user_id bigint NOT NULL,
  constraint fk_user
    foreign key (user_id)
    references users(id)
    on delete cascade,
  PRIMARY KEY(word_id, user_id),
  created_at timestamptz NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS review_words;
-- +goose StatementEnd
