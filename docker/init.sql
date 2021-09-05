-- user: tenant admin
CREATE USER "fims_tenant_manager" PASSWORD 'Welcome1!';
-- user: local tenant
CREATE USER "93fe5029-1886-4b63-94ca-35c503a52eff" PASSWORD 'c357e151-ff10-4603-aca3-d4b8f5ee676d';

-- create schema for tenant admin
CREATE SCHEMA "fims_tenant_manager" AUTHORIZATION "fims_tenant_manager";
-- create schema for user local tenant
CREATE SCHEMA "93fe5029-1886-4b63-94ca-35c503a52eff" AUTHORIZATION "93fe5029-1886-4b63-94ca-35c503a52eff";

-- revoke authorization for tenant admin
REVOKE ALL ON SCHEMA "fims_tenant_manager" FROM public;
-- revoke authorization for local tenant
REVOKE ALL ON SCHEMA "93fe5029-1886-4b63-94ca-35c503a52eff" FROM public;


-- create tenants table
CREATE TABLE "fims_tenant_manager"."tenants" (
    "id" uuid,
    "subdomain" text,
    "db_conn_password" text,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    PRIMARY KEY ("id")
);
-- init data for local tenant
INSERT INTO "fims_tenant_manager"."tenants" (
    "id",
    "subdomain",
    "db_conn_password",
    "created_at",
    "updated_at"
) VALUES (
    '93fe5029-1886-4b63-94ca-35c503a52eff',
    'localhost',
    'c357e151-ff10-4603-aca3-d4b8f5ee676d',
    '2021-08-29 15:16:58.159',
    '2021-08-29 15:16:58.159'
);

-- change owner
ALTER TABLE "fims_tenant_manager"."tenants" OWNER TO "fims_tenant_manager";