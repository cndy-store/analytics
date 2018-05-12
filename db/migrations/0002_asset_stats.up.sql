CREATE TABLE asset_stats (
  id SERIAL PRIMARY KEY,
  paging_token VARCHAR(128),
  asset_type VARCHAR(128),
  asset_code VARCHAR(128),
  asset_issuer VARCHAR(128),
  total_amount REAL,
  num_accounts INTEGER,
  num_effects INTEGER,
  created_at TIMESTAMP
);
