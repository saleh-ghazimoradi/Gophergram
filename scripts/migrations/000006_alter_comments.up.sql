ALTER TABLE comments
ADD CONSTRAINT fk_comments_user FOREIGN KEY (user_id)
REFERENCES users (id) ON DELETE CASCADE;
