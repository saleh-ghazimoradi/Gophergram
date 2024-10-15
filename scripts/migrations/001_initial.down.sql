DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.posts;
ALTER TABLE public.posts DROP CONSTRAINT fk_user;
ALTER TABLE public.posts DROP COLUMN tags;
Alter TABLE public.posts DROP COLUMN updated_at;
DROP TABLE IF EXISTS public.comments;
ALTER TABLE public.posts DROP COLUMN version;
DROP TABLE IF EXISTS followers;

CREATE INDEX IF EXISTS idx_posts_title;
CREATE INDEX IF EXISTS idx_posts_tags;
CREATE INDEX IF EXISTS idx_users_username;
CREATE INDEX IF EXISTS idx_posts_user_id;
CREATE INDEX IF EXISTS idx_comments_post_id;