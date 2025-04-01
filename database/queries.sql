--! Below are specific to authentication with AuthJS and NeonAdapter in NextJS

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--
-- Create Verification Token table (Tweets and Comments)
--
CREATE TABLE verification_token
(
  identifier TEXT NOT NULL,
  expires TIMESTAMPTZ NOT NULL,
  token TEXT NOT NULL,
 
  PRIMARY KEY (identifier, token)
);

--
-- Create Accounts table
--
CREATE TABLE accounts
(
  id SERIAL,
  "userId" UUID NOT NULL,
  type VARCHAR(255) NOT NULL,
  provider VARCHAR(255) NOT NULL,
  "providerAccountId" VARCHAR(255) NOT NULL,
  refresh_token TEXT,
  access_token TEXT,
  expires_at BIGINT,
  id_token TEXT,
  scope TEXT,
  session_state TEXT,
  token_type TEXT,
 
  PRIMARY KEY (id)
);

--
-- Create Sessions table
--
CREATE TABLE sessions
(
  id SERIAL,
  "userId" UUID NOT NULL,
  expires TIMESTAMPTZ NOT NULL,
  "sessionToken" VARCHAR(255) NOT NULL,
 
  PRIMARY KEY (id)
);

--
-- Create Users table
--
CREATE TABLE users
(
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  name VARCHAR(255),
  email VARCHAR(255),
  "emailVerified" TIMESTAMPTZ,
  image TEXT
);

--! Below are project specific

--
-- Alter the users table to add username and bio columns
--
ALTER TABLE users 
ADD COLUMN username TEXT UNIQUE, 
ADD COLUMN bio TEXT;

--
-- Create Posts table (Tweets and Comments)
--
CREATE TABLE IF NOT EXISTS posts (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid NOT NULL REFERENCES users(id),
  content text,
  parent_id uuid REFERENCES posts(id) ON DELETE CASCADE,
  likes_count integer DEFAULT 0,
  retweets_count integer DEFAULT 0,
  comments_count integer DEFAULT 0,
  created_at timestamp with time zone DEFAULT now()
);

ALTER TABLE posts 
ALTER COLUMN content SET NOT NULL;

--
-- Index for quickly fetching comments
--
CREATE INDEX idx_posts_parent_id ON posts(parent_id);

--
-- Create Likes table
--
CREATE TABLE IF NOT EXISTS likes (
  user_id uuid NOT NULL REFERENCES users(id),
  post_id uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  created_at timestamp with time zone DEFAULT now(),
  PRIMARY KEY (user_id, post_id)
);

--
-- Create Retweets table
--
CREATE TABLE IF NOT EXISTS retweets (
  user_id uuid NOT NULL REFERENCES users(id),
  post_id uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  created_at timestamp with time zone DEFAULT now(),
  PRIMARY KEY (user_id, post_id)
);

--
-- Create Follows table
--
CREATE TABLE IF NOT EXISTS follows (
  follower_id uuid NOT NULL REFERENCES users(id),
  following_id uuid NOT NULL REFERENCES users(id),
  created_at timestamp with time zone DEFAULT now(),
  PRIMARY KEY (follower_id, following_id)
);

--
-- Trigger to update post's likes count
--
CREATE OR REPLACE FUNCTION update_likes_count() RETURNS trigger AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    UPDATE posts SET likes_count = likes_count + 1 WHERE id = NEW.post_id;
  ELSIF TG_OP = 'DELETE' THEN
    UPDATE posts 
    SET likes_count = GREATEST(0, likes_count - 1) 
    WHERE id = OLD.post_id;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_likes_count
AFTER INSERT OR DELETE ON likes
FOR EACH ROW EXECUTE FUNCTION update_likes_count();

--
-- Trigger to update post's retweets count
--
CREATE OR REPLACE FUNCTION update_retweets_count() RETURNS trigger AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    UPDATE posts SET retweets_count = retweets_count + 1 WHERE id = NEW.post_id;
  ELSIF TG_OP = 'DELETE' THEN
    UPDATE posts 
    SET retweets_count = GREATEST(0, retweets_count - 1) 
    WHERE id = OLD.post_id;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_retweets_count
AFTER INSERT OR DELETE ON retweets
FOR EACH ROW EXECUTE FUNCTION update_retweets_count();

--
-- Trigger to update comments count
--
CREATE OR REPLACE FUNCTION update_comments_count() RETURNS trigger AS $$
BEGIN
  IF TG_OP = 'INSERT' AND NEW.parent_id IS NOT NULL THEN
    UPDATE posts SET comments_count = comments_count + 1 WHERE id = NEW.parent_id;
  ELSIF TG_OP = 'DELETE' AND OLD.parent_id IS NOT NULL THEN
    UPDATE posts 
    SET comments_count = GREATEST(0, comments_count - 1) 
    WHERE id = OLD.parent_id;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_comments_count
AFTER INSERT OR DELETE ON posts
FOR EACH ROW EXECUTE FUNCTION update_comments_count();

--
-- Alter Users Table to Add Follower/Following Count
--
ALTER TABLE users
ADD COLUMN IF NOT EXISTS follower_count integer DEFAULT 0,
ADD COLUMN IF NOT EXISTS following_count integer DEFAULT 0;

--
-- Trigger to update follower/following count
--
CREATE OR REPLACE FUNCTION update_follow_counts() RETURNS trigger AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    UPDATE users SET follower_count = follower_count + 1 WHERE id = NEW.following_id;
    UPDATE users SET following_count = following_count + 1 WHERE id = NEW.follower_id;
  ELSIF TG_OP = 'DELETE' THEN
    UPDATE users 
    SET follower_count = GREATEST(0, follower_count - 1) 
    WHERE id = OLD.following_id;
    UPDATE users 
    SET following_count = GREATEST(0, following_count - 1) 
    WHERE id = OLD.follower_id;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_follow_counts
AFTER INSERT OR DELETE ON follows
FOR EACH ROW EXECUTE FUNCTION update_follow_counts();
