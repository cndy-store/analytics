-- Type specifications are taken from
-- https://github.com/stellar/stellar-core/blob/master/docs/db-schema.md

CREATE TABLE effects (
  -- TODO: Use serial primary key?
  effect_id character varying(56) PRIMARY KEY,
  operation character varying(128),        -- TODO: Verify type with Stellar documentation
  succeeds character varying(128),         -- TODO: Verify type with Stellar documentation
  precedes character varying(128),         -- TODO: Verify type with Stellar documentation
  paging_token character varying(64),      -- TODO: Verify type with Stellar documentation
  account character varying(56),
  amount bigint,
  type character varying(128),             -- TODO: Verify type with Stellar documentation
  type_i integer,                          -- TODO: Verify type with Stellar documentation
  starting_balance character varying(128), -- TODO: Verify type with Stellar documentation

  balance BIGINT,
  balance_limit BIGINT,                    -- TODO: Verify type with Stellar documentation

  -- Types are taken from
  -- github.com/stellar/go/services/horizon/internal/db2/schema/migrations/6_create_assets_table.sql
  asset_type character varying(64),
  asset_code character varying(12),
  asset_issuer character varying(56),

  signer_public_key character varying(56),
  signer_weight integer,
  signer_key character varying(128),       -- TODO: Verify type with Stellar documentation
  signer_type character varying(128),      -- TODO: Verify type with Stellar documentation

  -- This field is extracted from the corresponding operation
  created_at timestamp without time zone
);
