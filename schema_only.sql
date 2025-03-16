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
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

-- *not* creating schema, since initdb creates it


--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON SCHEMA public IS '';


--
-- Name: user_role; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.user_role AS ENUM (
    'Farmer',
    'Factory',
    'Retailer',
    'Logistics',
    'Admin'
);


--
-- Name: create_yearly_sequence(text); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.create_yearly_sequence(year_prefix text) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
	EXECUTE 'CREATE SEQUENCE IF NOT EXISTS user_seq_' || year_prefix ||
        	' START 1;';
END;
$$;


--
-- Name: generate_userid(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.generate_userid() RETURNS text
    LANGUAGE plpgsql
    AS $$
DECLARE
	prefix TEXT;
	seq_num INTEGER;
	result TEXT;
BEGIN
	-- กำหนดคำนำหน้าด้วยปีปัจจุบันสองหลักสุดท้าย
	prefix := TO_CHAR(CURRENT_DATE, 'YY');

	-- ดึงค่าลำดับถัดไปจากลำดับที่กำหนดสำหรับปีปัจจุบัน
	seq_num := nextval('user_seq_' || prefix);

	-- สร้าง userid โดยรวมคำนำหน้าและหมายเลขลำดับ
	result := prefix || TO_CHAR(seq_num, 'FM0000');

	RETURN result;
END;
$$;


--
-- Name: generate_userid(text); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.generate_userid(role_prefix text) RETURNS text
    LANGUAGE plpgsql
    AS $$
DECLARE
	new_id INT;
	new_userid TEXT;
BEGIN
	SELECT NEXTVAL('user_seq') INTO new_id;
	new_userid := role_prefix || TO_CHAR(new_id, 'FM000');
	RETURN new_userid;
END;
$$;


--
-- Name: update_timestamp(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.categories (
    category_id bigint NOT NULL,
    name text NOT NULL
);


--
-- Name: categories_category_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.categories_category_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: categories_category_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.categories_category_id_seq OWNED BY public.categories.category_id;


--
-- Name: factory_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.factory_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: dairyfactory; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.dairyfactory (
    factoryid character varying(255) DEFAULT nextval('public.factory_id_seq'::regclass) NOT NULL,
    companyname text NOT NULL,
    address text NOT NULL,
    district text,
    subdistrict text,
    province text,
    country character varying(255) DEFAULT 'Thailand'::character varying,
    postcode text,
    telephone text NOT NULL,
    lineid text,
    facebook text,
    location_link text,
    createdon timestamp with time zone,
    wallet_address text NOT NULL,
    email text NOT NULL
);


--
-- Name: dairyfactory_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.dairyfactory_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: externalid; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.externalid (
    externalid character varying(255) NOT NULL,
    factoryid character varying(255),
    logisticname character varying(255),
    sendername character varying(255),
    logisticshippingdate date,
    logisticdeliverydate date,
    logisticqualitycheck boolean,
    logistictemp double precision,
    retailersreceiptdate date,
    retailerqualitycheck boolean,
    retailertemp double precision,
    retailername character varying(255),
    createdon timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: farmer_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.farmer_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: farmer; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.farmer (
    farmerid character varying(255) DEFAULT nextval('public.farmer_id_seq'::regclass) NOT NULL,
    companyname text NOT NULL,
    address text NOT NULL,
    district text,
    subdistrict text,
    province text,
    country character varying(255) DEFAULT 'Thailand'::character varying,
    postcode text,
    telephone text NOT NULL,
    lineid text,
    facebook text,
    location_link text,
    createdon timestamp with time zone,
    wallet_address text NOT NULL,
    email text NOT NULL
);


--
-- Name: logistics_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.logistics_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: logisticsprovider; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.logisticsprovider (
    logisticsid character varying(255) DEFAULT nextval('public.logistics_id_seq'::regclass) NOT NULL,
    companyname text NOT NULL,
    address text NOT NULL,
    district text,
    subdistrict text,
    province text,
    country character varying(255) DEFAULT 'Thailand'::character varying,
    postcode text,
    telephone text NOT NULL,
    lineid text,
    facebook text,
    location_link text,
    createdon timestamp with time zone,
    wallet_address text NOT NULL,
    email text NOT NULL
);


--
-- Name: organiccertification; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.organiccertification (
    certificationid character varying(255) NOT NULL,
    certificationtype text,
    certificationcid text,
    effective_date timestamp with time zone,
    issued_date timestamp with time zone,
    created_on timestamp with time zone,
    entityid text,
    entitytype text,
    blockchain_tx text
);


--
-- Name: product_lot_images; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_lot_images (
    id bigint NOT NULL,
    lot_id text NOT NULL,
    image_c_id text NOT NULL,
    created_at timestamp with time zone,
    tracking_ids text NOT NULL,
    person_in_charge text NOT NULL
);


--
-- Name: product_lot_images_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.product_lot_images_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: product_lot_images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.product_lot_images_id_seq OWNED BY public.product_lot_images.id;


--
-- Name: retailer_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.retailer_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: retailer; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.retailer (
    retailerid character varying(255) DEFAULT nextval('public.retailer_id_seq'::regclass) NOT NULL,
    companyname text NOT NULL,
    address text NOT NULL,
    district text,
    subdistrict text,
    province text,
    country character varying(255) DEFAULT 'Thailand'::character varying,
    postcode text,
    telephone text NOT NULL,
    lineid text,
    facebook text,
    location_link text,
    createdon timestamp with time zone,
    wallet_address text NOT NULL,
    email text NOT NULL
);


--
-- Name: tracking_status; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tracking_status (
    id integer NOT NULL,
    tracking_id text NOT NULL,
    status integer DEFAULT 0 NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: tracking_status_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tracking_status_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tracking_status_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tracking_status_id_seq OWNED BY public.tracking_status.id;


--
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    userid character varying(255) DEFAULT gen_random_uuid() NOT NULL,
    username text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    role text NOT NULL,
    entityid text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp without time zone,
    telephone text,
    profile_image bytea
);


--
-- Name: categories category_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories ALTER COLUMN category_id SET DEFAULT nextval('public.categories_category_id_seq'::regclass);


--
-- Name: product_lot_images id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_lot_images ALTER COLUMN id SET DEFAULT nextval('public.product_lot_images_id_seq'::regclass);


--
-- Name: tracking_status id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tracking_status ALTER COLUMN id SET DEFAULT nextval('public.tracking_status_id_seq'::regclass);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (category_id);


--
-- Name: dairyfactory dairyfactory_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.dairyfactory
    ADD CONSTRAINT dairyfactory_email_key UNIQUE (email);


--
-- Name: dairyfactory dairyfactory_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.dairyfactory
    ADD CONSTRAINT dairyfactory_pkey PRIMARY KEY (factoryid);


--
-- Name: externalid externalid_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.externalid
    ADD CONSTRAINT externalid_pkey PRIMARY KEY (externalid);


--
-- Name: farmer farmer_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.farmer
    ADD CONSTRAINT farmer_email_key UNIQUE (email);


--
-- Name: farmer farmer_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.farmer
    ADD CONSTRAINT farmer_pkey PRIMARY KEY (farmerid);


--
-- Name: logisticsprovider logisticsprovider_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logisticsprovider
    ADD CONSTRAINT logisticsprovider_email_key UNIQUE (email);


--
-- Name: logisticsprovider logisticsprovider_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logisticsprovider
    ADD CONSTRAINT logisticsprovider_pkey PRIMARY KEY (logisticsid);


--
-- Name: organiccertification organiccertification_certificationcid_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.organiccertification
    ADD CONSTRAINT organiccertification_certificationcid_key UNIQUE (certificationcid);


--
-- Name: organiccertification organiccertification_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.organiccertification
    ADD CONSTRAINT organiccertification_pkey PRIMARY KEY (certificationid);


--
-- Name: retailer retailer_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.retailer
    ADD CONSTRAINT retailer_email_key UNIQUE (email);


--
-- Name: retailer retailer_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.retailer
    ADD CONSTRAINT retailer_pkey PRIMARY KEY (retailerid);


--
-- Name: tracking_status tracking_status_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tracking_status
    ADD CONSTRAINT tracking_status_pkey PRIMARY KEY (id);


--
-- Name: tracking_status tracking_status_tracking_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tracking_status
    ADD CONSTRAINT tracking_status_tracking_id_key UNIQUE (tracking_id);


--
-- Name: categories uni_categories_name; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT uni_categories_name UNIQUE (name);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (userid);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: users set_timestamp; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_timestamp();


--
-- PostgreSQL database dump complete
--

