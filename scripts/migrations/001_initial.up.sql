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