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
-- Name: user_role; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.user_role AS ENUM (
    'Farmer',
    'Factory',
    'Retailer',
    'Logistics',
    'Admin'
);


ALTER TYPE public.user_role OWNER TO postgres;

--
-- Name: create_yearly_sequence(text); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.create_yearly_sequence(year_prefix text) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
	EXECUTE 'CREATE SEQUENCE IF NOT EXISTS user_seq_' || year_prefix ||
        	' START 1;';
END;
$$;


ALTER FUNCTION public.create_yearly_sequence(year_prefix text) OWNER TO postgres;

--
-- Name: generate_userid(); Type: FUNCTION; Schema: public; Owner: postgres
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


ALTER FUNCTION public.generate_userid() OWNER TO postgres;

--
-- Name: generate_userid(text); Type: FUNCTION; Schema: public; Owner: postgres
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


ALTER FUNCTION public.generate_userid(role_prefix text) OWNER TO postgres;

--
-- Name: update_timestamp(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_timestamp() OWNER TO postgres;

--
-- Name: dairyfactory_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.dairyfactory_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.dairyfactory_id_seq OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: dairyfactory; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.dairyfactory (
    factoryid character varying(255) DEFAULT nextval('public.dairyfactory_id_seq'::regclass) NOT NULL,
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


ALTER TABLE public.dairyfactory OWNER TO postgres;

--
-- Name: externalid; Type: TABLE; Schema: public; Owner: postgres
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


ALTER TABLE public.externalid OWNER TO postgres;

--
-- Name: farmer_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.farmer_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.farmer_id_seq OWNER TO postgres;

--
-- Name: farmer; Type: TABLE; Schema: public; Owner: postgres
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
    email text NOT NULL,
    entityid text NOT NULL
);


ALTER TABLE public.farmer OWNER TO postgres;

--
-- Name: logistics_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.logistics_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.logistics_id_seq OWNER TO postgres;

--
-- Name: logisticsprovider; Type: TABLE; Schema: public; Owner: postgres
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


ALTER TABLE public.logisticsprovider OWNER TO postgres;

--
-- Name: organiccertification; Type: TABLE; Schema: public; Owner: postgres
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


ALTER TABLE public.organiccertification OWNER TO postgres;

--
-- Name: retailer_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.retailer_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.retailer_id_seq OWNER TO postgres;

--
-- Name: retailer; Type: TABLE; Schema: public; Owner: postgres
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


ALTER TABLE public.retailer OWNER TO postgres;

--
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_id_seq OWNER TO postgres;

--
-- Name: user_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_seq OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
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
    deleted_at timestamp without time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Data for Name: dairyfactory; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.dairyfactory (factoryid, companyname, address, district, subdistrict, province, country, postcode, telephone, lineid, facebook, location_link, createdon, wallet_address, email) FROM stdin;
\.


--
-- Data for Name: externalid; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.externalid (externalid, factoryid, logisticname, sendername, logisticshippingdate, logisticdeliverydate, logisticqualitycheck, logistictemp, retailersreceiptdate, retailerqualitycheck, retailertemp, retailername, createdon) FROM stdin;
\.


--
-- Data for Name: farmer; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.farmer (farmerid, companyname, address, district, subdistrict, province, country, postcode, telephone, lineid, facebook, location_link, createdon, wallet_address, email, entityid) FROM stdin;
\.


--
-- Data for Name: logisticsprovider; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.logisticsprovider (logisticsid, companyname, address, district, subdistrict, province, country, postcode, telephone, lineid, facebook, location_link, createdon, wallet_address, email) FROM stdin;
\.


--
-- Data for Name: organiccertification; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.organiccertification (certificationid, certificationtype, certificationcid, effective_date, issued_date, created_on, entityid, entitytype, blockchain_tx) FROM stdin;
\.


--
-- Data for Name: retailer; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.retailer (retailerid, companyname, address, district, subdistrict, province, country, postcode, telephone, lineid, facebook, location_link, createdon, wallet_address, email) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (userid, username, email, password, role, entityid, created_at, updated_at, deleted_at) FROM stdin;
2500048	Apiwat Bunyasartpan	farm@test.com	$2a$14$qskTlxk/N1wPXaNIqeGFZ.2PZLKQqZ7GoSt9bH6FzmELptQt3Z0w6	pending	PENDING_ROLE	2025-02-26 13:23:50.617784+07	2025-02-26 13:23:50.617784+07	0001-01-01 00:00:00
\.


--
-- Name: dairyfactory_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.dairyfactory_id_seq', 24, true);


--
-- Name: farmer_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.farmer_id_seq', 59, true);


--
-- Name: logistics_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.logistics_id_seq', 4, true);


--
-- Name: retailer_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.retailer_id_seq', 3, true);


--
-- Name: user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_id_seq', 48, true);


--
-- Name: user_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_seq', 1, false);


--
-- Name: dairyfactory dairyfactory_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dairyfactory
    ADD CONSTRAINT dairyfactory_email_key UNIQUE (email);


--
-- Name: dairyfactory dairyfactory_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dairyfactory
    ADD CONSTRAINT dairyfactory_pkey PRIMARY KEY (factoryid);


--
-- Name: externalid externalid_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.externalid
    ADD CONSTRAINT externalid_pkey PRIMARY KEY (externalid);


--
-- Name: farmer farmer_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.farmer
    ADD CONSTRAINT farmer_email_key UNIQUE (email);


--
-- Name: farmer farmer_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.farmer
    ADD CONSTRAINT farmer_pkey PRIMARY KEY (farmerid);


--
-- Name: logisticsprovider logisticsprovider_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.logisticsprovider
    ADD CONSTRAINT logisticsprovider_email_key UNIQUE (email);


--
-- Name: logisticsprovider logisticsprovider_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.logisticsprovider
    ADD CONSTRAINT logisticsprovider_pkey PRIMARY KEY (logisticsid);


--
-- Name: organiccertification organiccertification_certificationcid_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.organiccertification
    ADD CONSTRAINT organiccertification_certificationcid_key UNIQUE (certificationcid);


--
-- Name: organiccertification organiccertification_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.organiccertification
    ADD CONSTRAINT organiccertification_pkey PRIMARY KEY (certificationid);


--
-- Name: retailer retailer_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.retailer
    ADD CONSTRAINT retailer_email_key UNIQUE (email);


--
-- Name: retailer retailer_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.retailer
    ADD CONSTRAINT retailer_pkey PRIMARY KEY (retailerid);


--
-- Name: farmer uni_farmer_entityid; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.farmer
    ADD CONSTRAINT uni_farmer_entityid UNIQUE (entityid);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (userid);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: users set_timestamp; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_timestamp();


--
-- PostgreSQL database dump complete
--
