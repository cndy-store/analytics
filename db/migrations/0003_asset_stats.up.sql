CREATE TABLE asset_stats (
  paging_token character varying(64) PRIMARY KEY,
  asset_type character varying(64),
  asset_code character varying(12),
  asset_issuer character varying(56),
  total_amount bigint,
  num_accounts integer,
  num_effects integer,
  created_at timestamp without time zone
);
