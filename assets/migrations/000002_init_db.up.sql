-- This is required to be able to run the function gen_random_bytes
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- https://stackoverflow.com/questions/41970461/how-to-generate-a-random-unique-alphanumeric-id-of-length-n-in-postgres-9-6#41988979
CREATE
OR REPLACE FUNCTION generate_uid(size INT) RETURNS TEXT AS
$$
DECLARE
characters TEXT  := 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    bytes
BYTEA := gen_random_bytes(size);
    l
INT   := length(characters);
    i
INT   := 0;
output     TEXT  := '';
BEGIN
    WHILE
i < size
        LOOP
            output := output || substr(characters, get_byte(bytes, i) % l + 1, 1);
            i
:= i + 1;
END LOOP;
    RETURN
output;
END;
$$
LANGUAGE plpgsql VOLATILE;

-- https://x-team.com/blog/automatic-timestamps-with-postgresql/
CREATE
OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.UPDATED_AT = NOW();
RETURN NEW;
END;
$$
LANGUAGE plpgsql;
