CREATE TABLE IF NOT EXISTS posts (
  id bigserial PRIMARY KEY,
  content TEXT NOT NULL,
    title TEXT NOT NULL,
    user_id bigint NOT NULL
);

ALTER TABLE posts ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id);