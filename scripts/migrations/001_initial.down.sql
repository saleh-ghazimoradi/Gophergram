DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.posts;
ALTER TABLE public.posts DROP CONSTRAINT fk_user;
ALTER TABLE public.posts DROP COLUMN tags;
Alter TABLE public.posts DROP COLUMN updated_at;