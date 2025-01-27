--
-- PostgreSQL database dump
--

-- Dumped from database version 16.6 (Ubuntu 16.6-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.6 (Ubuntu 16.6-0ubuntu0.24.04.1)

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
-- Name: order_status; Type: TYPE; Schema: public; Owner: gratify
--

CREATE TYPE public.order_status AS ENUM (
    'NEW',
    'REGISTERED',
    'PROCESSING',
    'PROCESSED',
    'INVALID'
);


ALTER TYPE public.order_status OWNER TO gratify;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: accounts; Type: TABLE; Schema: public; Owner: gratify
--

CREATE TABLE public.accounts (
    id integer NOT NULL,
    user_id bigint NOT NULL,
    amount money DEFAULT '$0.00'::money NOT NULL
);


ALTER TABLE public.accounts OWNER TO gratify;

--
-- Name: accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: gratify
--

CREATE SEQUENCE public.accounts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.accounts_id_seq OWNER TO gratify;

--
-- Name: accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: gratify
--

ALTER SEQUENCE public.accounts_id_seq OWNED BY public.accounts.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: gratify
--

CREATE TABLE public.orders (
    id bigint NOT NULL,
    user_id integer NOT NULL,
    status public.order_status DEFAULT 'NEW'::public.order_status NOT NULL,
    accrual money DEFAULT '$0.00'::money NOT NULL,
    sum money DEFAULT '$0.00'::money NOT NULL,
    date timestamp with time zone NOT NULL
);


ALTER TABLE public.orders OWNER TO gratify;

--
-- Name: orders_user_id_seq; Type: SEQUENCE; Schema: public; Owner: gratify
--

CREATE SEQUENCE public.orders_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.orders_user_id_seq OWNER TO gratify;

--
-- Name: orders_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: gratify
--

ALTER SEQUENCE public.orders_user_id_seq OWNED BY public.orders.user_id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: gratify
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO gratify;

--
-- Name: users; Type: TABLE; Schema: public; Owner: gratify
--

CREATE TABLE public.users (
    id integer NOT NULL,
    login character varying(64) NOT NULL,
    password character varying(64) NOT NULL
);


ALTER TABLE public.users OWNER TO gratify;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: gratify
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO gratify;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: gratify
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: accounts id; Type: DEFAULT; Schema: public; Owner: gratify
--

ALTER TABLE ONLY public.accounts ALTER COLUMN id SET DEFAULT nextval('public.accounts_id_seq'::regclass);


--
-- Name: orders user_id; Type: DEFAULT; Schema: public; Owner: gratify
--

ALTER TABLE ONLY public.orders ALTER COLUMN user_id SET DEFAULT nextval('public.orders_user_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: gratify
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: gratify
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: accounts accounts_user_id_key; Type: CONSTRAINT; Schema: public; Owner: gratify
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_user_id_key UNIQUE (user_id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: gratify
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: gratify
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: users users_login_key; Type: CONSTRAINT; Schema: public; Owner: gratify
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_login_key UNIQUE (login);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: gratify
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: accounts_user_id; Type: INDEX; Schema: public; Owner: gratify
--

CREATE INDEX accounts_user_id ON public.accounts USING btree (user_id);


--
-- Name: orders_date; Type: INDEX; Schema: public; Owner: gratify
--

CREATE INDEX orders_date ON public.orders USING btree (date);


--
-- Name: orders_id; Type: INDEX; Schema: public; Owner: gratify
--

CREATE INDEX orders_id ON public.orders USING btree (id);


--
-- Name: orders_user_id; Type: INDEX; Schema: public; Owner: gratify
--

CREATE INDEX orders_user_id ON public.orders USING btree (user_id);


--
-- Name: users_id; Type: INDEX; Schema: public; Owner: gratify
--

CREATE INDEX users_id ON public.users USING btree (id);


--
-- Name: users_login; Type: INDEX; Schema: public; Owner: gratify
--

CREATE INDEX users_login ON public.users USING btree (login);


--
-- Name: accounts accounts_fk1; Type: FK CONSTRAINT; Schema: public; Owner: gratify
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_fk1 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: orders orders_fk0; Type: FK CONSTRAINT; Schema: public; Owner: gratify
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_fk0 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

