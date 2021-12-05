-- user: fims admin
CREATE USER "fims" PASSWORD 'DtrEmef4Q4avBmUK';
-- user: local tenant
CREATE USER "93fe5029-1886-4b63-94ca-35c503a52eff" PASSWORD 'hePrAqafu&5Ep49V8th9';

-- create schema for tenant admin
CREATE SCHEMA "fims" AUTHORIZATION "fims";
-- create schema for user local tenant
CREATE SCHEMA "93fe5029-1886-4b63-94ca-35c503a52eff" AUTHORIZATION "93fe5029-1886-4b63-94ca-35c503a52eff";

-- revoke authorization for tenant admin
REVOKE ALL ON SCHEMA "fims" FROM public;
-- revoke authorization for local tenant
REVOKE ALL ON SCHEMA "93fe5029-1886-4b63-94ca-35c503a52eff" FROM public;


-- create tenants table
CREATE TABLE "fims"."tenants" (
    "id" uuid,
    "subdomain" text,
    "dsn" text,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    PRIMARY KEY ("id")
);
-- init data for local tenant
INSERT INTO "fims"."tenants" (
    "id",
    "subdomain",
    "dsn",
    "created_at",
    "updated_at"
) VALUES (
    '93fe5029-1886-4b63-94ca-35c503a52eff',
    'localhost',
    'host=localhost port=5432 user=93fe5029-1886-4b63-94ca-35c503a52eff password=hePrAqafu&5Ep49V8th9 dbname=postgres sslmode=disable TimeZone=UTC',
    '2021-08-29 15:16:58.159',
    '2021-08-29 15:16:58.159'
);

-- change owner
ALTER TABLE "fims"."tenants" OWNER TO "fims";