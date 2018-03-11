CREATE TABLE asset_stats (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  asset_type VARCHAR(128),
  asset_code VARCHAR(128),
  asset_issuer VARCHAR(128),
  total_amount VARCHAR(128),
  num_accounts INTEGER,
  num_effects INTEGER,
  created_at DATETIME
);
