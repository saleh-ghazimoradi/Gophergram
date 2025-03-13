CREATE TABLE IF NOT EXISTS user_invitations (
    token TEXT PRIMARY KEY,
    user_id bigint NOT NULL
);