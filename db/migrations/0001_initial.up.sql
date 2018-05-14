CREATE TABLE cursors (
  id SERIAL PRIMARY KEY,
  paging_token VARCHAR(64) NOT NULL
);

-- Set genesis cursor (Mon Mar 12 18:01:12 CET 2018)
INSERT INTO cursors(paging_token) VALUES('33819440072110101-2');

CREATE TABLE effects (
  id SERIAL PRIMARY KEY,
  effect_id VARCHAR(128),
  operation VARCHAR(128),
  succeeds VARCHAR(128),
  precedes VARCHAR(128),
  paging_token VARCHAR(128),
  account VARCHAR(128),
  amount REAL,
  type VARCHAR(128),
  type_i INTEGER,
  starting_balance VARCHAR(128),

  balance VARCHAR(128),
  balance_limit VARCHAR(128),

  asset_type VARCHAR(128),
  asset_code VARCHAR(128),
  asset_issuer VARCHAR(128),

  signer_public_key VARCHAR(128),
  signer_weight INTEGER,
  signer_key VARCHAR(128),
  signer_type VARCHAR(128),

  --- This field is extracted from the corresponding operation
  created_at TIMESTAMP
);
