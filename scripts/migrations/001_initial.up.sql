CREATE TABLE IF NOT EXISTS public.users (
    id bigserial PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    username varchar(255) UNIQUE NOT NULL,
    password bytea NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
    );

-- Create citext extension (for case-insensitive text)
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS public.posts (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    user_id bigint NOT NULL,
    content text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
    );

-- Add foreign key constraint for user_id referencing users table
ALTER TABLE public.posts ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users(id);

-- Add new columns to the posts table
ALTER TABLE public.posts ADD COLUMN tags VARCHAR(100) [];  -- Array of tags
ALTER TABLE public.posts ADD COLUMN updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW();  -- Timestamp for the last update

CREATE TABLE IF NOT EXISTS public.comments(
    id bigserial PRIMARY KEY,
    post_id bigserial NOT NULL,
    user_id bigserial NOT NULL,
    content TEXT NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

ALTER TABLE public.posts ADD COLUMN version INT DEFAULT 0;

CREATE TABLE IF NOT EXISTS public.followers (
    user_id bigint NOT NULL,
    follower_id bigint NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, follower_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_comments_content ON comments USING gin (content gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING gin (tags);

CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);

CREATE TABLE IF NOT EXISTS public.user_invitation (
    token bytea PRIMARY KEY,
    user_id bigint NOT NULL
);

ALTER TABLE public.users ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE public.user_invitation ADD COLUMN expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL;



CREATE TABLE IF NOT EXISTS public.roles (
  id BIGSERIAL PRIMARY KEY,
  name VACHAR(255) NOT NULL UNIQUE,
  level int NOT NULL DEFAULT 0,
  description TEXT
);

INSERT INTO public.roles (name, description, level) VALUES ('user', 'A user can create posts and comments', 1);
INSERT INTO public.roles (name, description, level) VALUES ('moderator', 'A moderator can update other users posts',2);
INSERT INTO public.roles (name, description, level) VALUES ('admin', 'An admin can update and delete other users posts', 3);

ALTER TABLE IF EXISTS public.users ADD COLUMN role_id INT REFERENCES roles(id) DEFAULT 1;

UPDATE public.users SET role_id = (
    SELECT id FROM roles WHERE name = 'user'
);


ALTER TABLE public.users ALTER COLUMN role_id DROP DEFAULT;

ALTER TABLE public.users ALTER COLUMN role_id SET NOT NULL;