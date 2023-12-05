CREATE EXTENSION IF NOT EXISTS dblink;

DO $$
BEGIN
PERFORM dblink_exec('', 'CREATE DATABASE testdb');
EXCEPTION WHEN duplicate_database THEN RAISE NOTICE '%, skipping', SQLERRM USING ERRCODE = SQLSTATE;
END
$$;

-- Connect to the 'sm_user_sb' database
\c testdb;

-- Create the 'sm_user2' user
CREATE USER sm_user3 WITH PASSWORD '12345678';

-- Grant CONNECT privilege to the user on the database
GRANT CONNECT ON DATABASE testdb TO sm_user3;

-- Grant USAGE privilege on the schema to the user
GRANT USAGE ON SCHEMA public TO sm_user3;

-- Grant CREATE privilege on the schema to the user
GRANT CREATE ON SCHEMA public TO sm_user3;

-- Grant SELECT, INSERT, UPDATE, DELETE privileges on all tables in the schema
GRANT  SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO sm_user3;

-- Optionally, grant USAGE privilege on sequences in the schema
GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO sm_user3;
