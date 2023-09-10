
SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: generate_uid(integer); Type: FUNCTION; Schema: public; Owner: rotabot
--

CREATE FUNCTION public.generate_uid(size integer) RETURNS text
    LANGUAGE plpgsql
    AS $$
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
$$;


ALTER FUNCTION public.generate_uid(size integer) OWNER TO rotabot;

--
-- Name: trigger_set_timestamp(); Type: FUNCTION; Schema: public; Owner: rotabot
--

CREATE FUNCTION public.trigger_set_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.UPDATED_AT = NOW();
RETURN NEW;
END;
$$;


ALTER FUNCTION public.trigger_set_timestamp() OWNER TO rotabot;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: rotas; Type: TABLE; Schema: public; Owner: rotabot
--

CREATE TABLE public.rotas (
    id text DEFAULT ('RT'::text || public.generate_uid(14)) NOT NULL,
    team_id text NOT NULL,
    channel_id text NOT NULL,
    name text NOT NULL,
    metadata jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.rotas OWNER TO rotabot;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: rotabot
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO rotabot;

--
-- Data for Name: rotas; Type: TABLE DATA; Schema: public; Owner: rotabot
--

COPY public.rotas (id, team_id, channel_id, name, metadata, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: rotabot
--

COPY public.schema_migrations (version, dirty) FROM stdin;
3	f
\.


--
-- Name: rotas rotas_pkey; Type: CONSTRAINT; Schema: public; Owner: rotabot
--

ALTER TABLE ONLY public.rotas
    ADD CONSTRAINT rotas_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: rotabot
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: idx_unique_rota_within_team_and_channel; Type: INDEX; Schema: public; Owner: rotabot
--

CREATE UNIQUE INDEX idx_unique_rota_within_team_and_channel ON public.rotas USING btree (name, channel_id, team_id);


--
-- Name: rotas rotas_updated_at_trigger; Type: TRIGGER; Schema: public; Owner: rotabot
--

CREATE TRIGGER rotas_updated_at_trigger BEFORE UPDATE ON public.rotas FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- PostgreSQL database dump complete
--

