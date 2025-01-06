CREATE TABLE IF NOT EXISTS user_invitations (
    token bytea PRIMARY KEY,
    user_id bigint NOT NULL
);

ALTER TABLE user_invitations ADD COLUMN expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL;