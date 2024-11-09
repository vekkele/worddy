-- +goose Up
CREATE TABLE IF NOT EXISTS stages (
  id bigserial PRIMARY KEY,
  level integer UNIQUE NOT NULL,
  hours_to_next integer NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS words (
  id bigserial PRIMARY KEY,
  word text NOT NULL,
  next_review timestamptz NOT NULL,
  stage_id bigint NOT NULL,
  constraint fk_stage
    foreign key (stage_id)
    references stages(id)
    on delete cascade,
  user_id bigint NOT NULL,
  constraint fk_user
    foreign key (user_id)
    references users(id)
    on delete cascade
);

CREATE TABLE IF NOT EXISTS translations (
  id bigserial PRIMARY KEY,
  translation text NOT NULL,
  word_id bigint NOT NULL,
  constraint fk_word
    foreign key (word_id)
    references words(id)
    on delete cascade
);

-- +goose statementbegin
DO $$
DECLARE day integer;
DECLARE week integer;
DECLARE month integer;
BEGIN
  day := 24;
  week := day * 7;
  month := day * 30;

  INSERT INTO stages (level, hours_to_next) VALUES (1, 4);
  INSERT INTO stages (level, hours_to_next) VALUES (2, 8);
  INSERT INTO stages (level, hours_to_next) VALUES (3, day);
  INSERT INTO stages (level, hours_to_next) VALUES (4, day * 2);
  INSERT INTO stages (level, hours_to_next) VALUES (5, week);
  INSERT INTO stages (level, hours_to_next) VALUES (6, week * 2);
  INSERT INTO stages (level, hours_to_next) VALUES (7, month * 1);
  INSERT INTO stages (level, hours_to_next) VALUES (8, month * 4);
  INSERT INTO stages (level, hours_to_next) VALUES (9, 0);
END $$;
-- +goose statementend


-- +goose Down
DROP TABLE IF EXISTS translations;
DROP TABLE IF EXISTS words;
DROP TABLE IF EXISTS stages;