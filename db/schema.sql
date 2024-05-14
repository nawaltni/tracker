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
-- Name: tiger; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA tiger;


--
-- Name: tiger_data; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA tiger_data;


--
-- Name: topology; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA topology;


--
-- Name: SCHEMA topology; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON SCHEMA topology IS 'PostGIS Topology schema';


--
-- Name: fuzzystrmatch; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS fuzzystrmatch WITH SCHEMA public;


--
-- Name: EXTENSION fuzzystrmatch; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION fuzzystrmatch IS 'determine similarities and distance between strings';


--
-- Name: postgis; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA public;


--
-- Name: EXTENSION postgis; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION postgis IS 'PostGIS geometry and geography spatial types and functions';


--
-- Name: postgis_tiger_geocoder; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS postgis_tiger_geocoder WITH SCHEMA tiger;


--
-- Name: EXTENSION postgis_tiger_geocoder; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION postgis_tiger_geocoder IS 'PostGIS tiger geocoder and reverse geocoder';


--
-- Name: postgis_topology; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS postgis_topology WITH SCHEMA topology;


--
-- Name: EXTENSION postgis_topology; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION postgis_topology IS 'PostGIS topology spatial types and functions';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: phone_meta; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.phone_meta (
    user_id uuid NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    id text NOT NULL,
    device_id text NOT NULL,
    brand text NOT NULL,
    model text NOT NULL,
    os text NOT NULL,
    app_version text NOT NULL,
    carrier text NOT NULL,
    battery integer NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(128) NOT NULL
);


--
-- Name: user_positions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_positions (
    user_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    reference text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    place_id uuid,
    place_name text,
    checked_in timestamp with time zone,
    checked_out timestamp with time zone,
    location public.geometry(Point,4326),
    backend_user_id character varying(255) NOT NULL
);


--
-- Name: phone_meta phone_meta_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.phone_meta
    ADD CONSTRAINT phone_meta_pkey PRIMARY KEY (user_id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: user_positions user_positions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_positions
    ADD CONSTRAINT user_positions_pkey PRIMARY KEY (user_id);


--
-- Name: idx_phone_meta_device_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_phone_meta_device_id ON public.phone_meta USING btree (device_id);


--
-- Name: idx_user_positions_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_positions_created_at ON public.user_positions USING btree (created_at);


--
-- Name: idx_user_positions_location; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_positions_location ON public.user_positions USING gist (location);


--
-- Name: idx_user_positions_reference; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_positions_reference ON public.user_positions USING btree (reference);


--
-- Name: user_positions_backend_user_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX user_positions_backend_user_id_index ON public.user_positions USING btree (backend_user_id);


--
-- Name: phone_meta phone_meta_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.phone_meta
    ADD CONSTRAINT phone_meta_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.user_positions(user_id) ON DELETE CASCADE;


--
-- Name: phone_meta phone_meta_user_id_fkey1; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.phone_meta
    ADD CONSTRAINT phone_meta_user_id_fkey1 FOREIGN KEY (user_id) REFERENCES public.user_positions(user_id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20231107030909'),
    ('20240402013652');
