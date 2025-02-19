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
    userid text NOT NULL,
    username text,
    companyname text,
    address text NOT NULL,
    city text,
    province text,
    country character varying(255) DEFAULT 'Thailand'::character varying NOT NULL,
    postcode text,
    email text,
    telephone text,
    lineid text,
    facebook text,
    location_link text,
    createdon timestamp with time zone,
    factory_id character varying(255) DEFAULT nextval('public.dairyfactory_id_seq'::regclass),
    company_name text,
    line_id text,
    created_on timestamp with time zone DEFAULT CURRENT_TIMESTAMP
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
    userid text NOT NULL,
    farmer_name text,
    companyname text,
    address text NOT NULL,
    city text,
    province text,
    country character varying(255) DEFAULT 'Thailand'::character varying NOT NULL,
    postcode text,
    email text,
    telephone text,
    lineid text,
    facebook text,
    location_link text,
    createdon timestamp with time zone,
    wallet_address text NOT NULL
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
    userid text,
    companyname text,
    telephone text,
    createdon timestamp with time zone,
    address text NOT NULL,
    city text,
    province text,
    country text DEFAULT 'Thailand'::text,
    postcode text,
    email text,
    lineid text,
    facebook text,
    location_link text,
    username text NOT NULL
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
    userid text,
    companyname text,
    telephone text,
    createdon timestamp with time zone,
    address text NOT NULL,
    city text,
    province text,
    country text DEFAULT 'Thailand'::text,
    postcode text,
    email text,
    lineid text,
    facebook text,
    location_link text,
    username text NOT NULL
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
    userid character varying(255) DEFAULT public.generate_userid() NOT NULL,
    email text,
    password text,
    role text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Data for Name: dairyfactory; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.dairyfactory (factoryid, userid, username, companyname, address, city, province, country, postcode, email, telephone, lineid, facebook, location_link, createdon, factory_id, company_name, line_id, created_on) FROM stdin;
\.


--
-- Data for Name: externalid; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.externalid (externalid, factoryid, logisticname, sendername, logisticshippingdate, logisticdeliverydate, logisticqualitycheck, logistictemp, retailersreceiptdate, retailerqualitycheck, retailertemp, retailername, createdon) FROM stdin;
\.


--
-- Data for Name: farmer; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.farmer (farmerid, userid, farmer_name, companyname, address, city, province, country, postcode, email, telephone, lineid, facebook, location_link, createdon, wallet_address) FROM stdin;
FA2500056	2500043	Apiwat Bunyasartpan	Mark Farm	92/395, lakhok muangeak	Kamphaeng Phet	Kamphaeng Phet	Thailand	12000	test@gmail.com	+66 0634371654	\N	\N	\N	2025-02-15 02:50:48.201208+07	0xCF74E3c1769bD376BB0CDf601DD3B081F2f8B0A6
\.


--
-- Data for Name: logisticsprovider; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.logisticsprovider (logisticsid, userid, companyname, telephone, createdon, address, city, province, country, postcode, email, lineid, facebook, location_link, username) FROM stdin;
\.


--
-- Data for Name: organiccertification; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.organiccertification (certificationid, certificationtype, certificationcid, effective_date, issued_date, created_on, entityid, entitytype, blockchain_tx) FROM stdin;
EVENT-1739303823	Organic	QmCertification123	2026-06-01 07:00:00+07	2024-06-01 07:00:00+07	2025-02-12 02:57:03.566626+07	FARM001	Farmer	0xb52b648fbb1d5156d7a33057f5a47ad5ce8ed05194df07348ed5672f188e7716
EVENT-1739304278	Organic	QmSz4an3jPbXre7Y3voKFierDiiXZ5M6eTUqy1SMgGpMjS	2026-06-01 07:00:00+07	2025-02-07 07:00:00+07	2025-02-12 03:04:39.031691+07	2500042	Farmer	0x559ed116202ad72f17f8c299c6243e66ff89bd8e2d0762d071a4091c60976cc5
EVENT-1739313952	Organic	QmWN1GPFzRAL7KbYqGyHBky1ve29Y2oBTD78Bqhzdpmhb4	2026-06-01 07:00:00+07	2025-02-07 07:00:00+07	2025-02-12 05:45:52.884607+07	FA2500043	Farmer	0x7f0da4417db9f2915b5a62c7ef61f4e9072d9c11b4ebc123bf1c2ba5206a89cb
EVENT-1739316187	Organic	QmTq37LUDjZWryBjMrtUXfdRffGcR3U2eT3WeyVh1eKeEm	2026-06-01 07:00:00+07	2025-02-07 07:00:00+07	2025-02-12 06:23:07.359519+07	DF2500023	Factory	0xa583520ca2a332574cd17def537de33e817ab0295f80ab1bbb64c6107fa884bf
EVENT-1739317467	Organic	QmdCJUzhqsJCbPUG4hBgopx1o5BDvMumriyy2iB4E27Ejv	2026-06-01 07:00:00+07	2025-02-07 07:00:00+07	2025-02-12 06:44:27.219699+07	RT2500003	Retailer	0x6a7c7bde57e148b3923358f73cce5ede42b7cea7ae2df17e9772249fea3f7f47
EVENT-1739319085	Organic	QmSE6s4sDT8giJFTkubWkNNxNovcqMBHWLmZyJhrisQRsi	2026-06-01 07:00:00+07	2025-02-07 07:00:00+07	2025-02-12 07:11:25.778886+07	LG2500001	Logistics	0xc5a3b4d389b1eb2b197b12ed7b8e62e7c1211f32200c43ff80dc42542c96f91a
EVENT-1739338141	Organic	QmYfvDGiBMDJLwG1xH4f4QzJWVhv2z2XkEZEdqRjwqon99	2026-06-01 07:00:00+07	2025-02-07 07:00:00+07	2025-02-12 12:29:01.757798+07	LG2500004	Logistics	0x3dd3d43aa72321edfc64f32ecd84e0372d297a581a399344d01fabd7d8e1c48e
EVENT-1739558620	Organic	Qme83eJ6RWwGd3R3mvKd8JcW2S5nhqoagCno2frvrK7sdj	2026-06-01 07:00:00+07	2025-02-07 07:00:00+07	2025-02-15 01:43:40.913789+07	FA2500055	Farmer	0x289e574dfb86a668295506d0704b4d36a4520e929b6af15100c0207364542825
\.


--
-- Data for Name: retailer; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.retailer (retailerid, userid, companyname, telephone, createdon, address, city, province, country, postcode, email, lineid, facebook, location_link, username) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (userid, email, password, role, created_at, updated_at, deleted_at) FROM stdin;
2500043	test@gmail.com	$2a$14$TIkY1keCg0iy2AtubF7xAuMud0SCx.pEVlrLoOwqgl3ppiWHIPCEq	farmer	2025-02-15 02:50:18.788822+07	2025-02-15 02:50:48.198377+07	0001-01-01 06:42:04+06:42:04
\.


--
-- Name: dairyfactory_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.dairyfactory_id_seq', 24, true);


--
-- Name: farmer_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.farmer_id_seq', 56, true);


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

SELECT pg_catalog.setval('public.user_id_seq', 43, true);


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
-- Name: dairyfactory dairyfactory_userid_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dairyfactory
    ADD CONSTRAINT dairyfactory_userid_key UNIQUE (userid);


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
-- Name: farmer farmer_userid_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.farmer
    ADD CONSTRAINT farmer_userid_key UNIQUE (userid);


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
-- Name: logisticsprovider uni_logisticsprovider_userid; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.logisticsprovider
    ADD CONSTRAINT uni_logisticsprovider_userid UNIQUE (userid);


--
-- Name: retailer uni_retailer_userid; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.retailer
    ADD CONSTRAINT uni_retailer_userid UNIQUE (userid);


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
-- Name: dairyfactory dairyfactory_userid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.dairyfactory
    ADD CONSTRAINT dairyfactory_userid_fkey FOREIGN KEY (userid) REFERENCES public.users(userid) ON DELETE CASCADE;


--
-- Name: externalid externalid_factoryid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.externalid
    ADD CONSTRAINT externalid_factoryid_fkey FOREIGN KEY (factoryid) REFERENCES public.dairyfactory(factoryid) ON DELETE CASCADE;


--
-- Name: farmer farmer_userid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.farmer
    ADD CONSTRAINT farmer_userid_fkey FOREIGN KEY (userid) REFERENCES public.users(userid) ON DELETE CASCADE;


--
-- Name: logisticsprovider logisticsprovider_userid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.logisticsprovider
    ADD CONSTRAINT logisticsprovider_userid_fkey FOREIGN KEY (userid) REFERENCES public.users(userid) ON DELETE CASCADE;


--
-- Name: retailer retailer_userid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.retailer
    ADD CONSTRAINT retailer_userid_fkey FOREIGN KEY (userid) REFERENCES public.users(userid) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

