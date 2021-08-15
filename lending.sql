CREATE DATABASE lending;

CREATE TABLE lending.public.account (
	account_id serial NOT NULL,
	first_name varchar(100) NOT NULL,
	last_name varchar(100) NOT NULL,
	phone varchar(20) NOT NULL,
	email varchar(100) NOT NULL,
	"password" varchar(200) NOT NULL,
	account_number varchar(30) NOT NULL,
	is_verify bool NOT NULL DEFAULT false,
	status varchar(20) NOT NULL DEFAULT 'PENDING'::character varying,
	CONSTRAINT account_email_key UNIQUE (email),
	CONSTRAINT account_id PRIMARY KEY (account_id)
);

CREATE TABLE lending.public.account_document (
	account_id int4 NOT NULL,
	document_id int4 NOT NULL,
	file_name varchar NOT NULL,
	file_context varchar NOT NULL,
	tag varchar NULL,
	CONSTRAINT account_document_pkey PRIMARY KEY (account_id, document_id)
);

CREATE TABLE lending.public.contract (
	contract_id serial NOT NULL,
	account_id int4 NOT NULL,
	interest_code int4 NOT NULL,
	loan_outstanding numeric NOT NULL,
	term int4 NOT NULL,
	status varchar(30) NOT NULL DEFAULT 'PENDING'::character varying,
	created_datetime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_datetime timestamp NULL,
	CONSTRAINT contract_pkey PRIMARY KEY (contract_id)
);

CREATE TABLE lending.public.document_info (
	document_id serial NOT NULL,
	document_type varchar(30) NOT NULL,
	CONSTRAINT document_info_pkey PRIMARY KEY (document_id)
);

CREATE TABLE lending.public.interest_term (
	interest_code serial NOT NULL,
	interest_rate numeric NOT NULL,
	CONSTRAINT interest_term_pkey PRIMARY KEY (interest_code)
);

CREATE TABLE lending.public.repay_transaction (
	id serial NOT NULL,
	contract_id int4 NOT NULL,
	account_id int4 NOT NULL,
	amount numeric NOT NULL,
	slip varchar NOT NULL,
	status varchar(30) NOT NULL DEFAULT 'PENDING'::character varying,
	created_datetime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_datetime timestamp NULL,
	CONSTRAINT repay_transaction_pkey PRIMARY KEY (id)
);

CREATE TABLE lending.public.terms_condition (
	account_id int4 NOT NULL,
	current_accept_version varchar(20) NOT NULL,
	CONSTRAINT term_condition_pkey PRIMARY KEY (account_id)
);

CREATE TABLE lending.public.wallet (
	account_id int4 NOT NULL,
	btc_volume numeric NOT NULL,
	eth_volume numeric NOT NULL,
	margin_call_date date NULL,
	latest_datetime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT wallet_pkey PRIMARY KEY (account_id)
);

CREATE TABLE lending.public.wallet_transaction (
	id serial NOT NULL,
	account_id int4 NOT NULL,
	address varchar(100) NOT NULL,
	chain_id int4 NOT NULL,
	txn_hash varchar(100) NULL,
	collateral_type varchar(10) NOT NULL,
	volume numeric NOT NULL,
	txn_type varchar(30) NOT NULL,
	status varchar(30) NOT NULL DEFAULT 'PENDING'::character varying,
	created_datetime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_datetime timestamp NULL,
	CONSTRAINT wallet_transaction_pkey PRIMARY KEY (id)
);

CREATE TABLE lending.public.user_subscription (
	first_name varchar(100) NOT NULL,
	last_name varchar(100) NOT NULL,
	phone varchar(20) NOT NULL,
	email varchar(100) NOT NULL,
	created_datetime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT user_subscription_pkey PRIMARY KEY (email)
);