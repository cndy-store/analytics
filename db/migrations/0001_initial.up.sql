CREATE TABLE cursors (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  paging_token VARCHAR(64) NOT NULL
);

CREATE TABLE effects (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  effect_id VARCHAR(128),
  operation VARCHAR(128),
  succeeds VARCHAR(128), -- Currently not used
  precedes VARCHAR(128), -- Currently not used
  paging_token VARCHAR(128),
  account VARCHAR(128),
  amount VARCHAR(128),
  type VARCHAR(128),
  type_i INTEGER, -- Currently not used
  starting_balance VARCHAR(128),

  balance VARCHAR(128),
  balance_limit VARCHAR(128),

  asset_type VARCHAR(128),
  asset_code VARCHAR(128),
  asset_issuer VARCHAR(128),

  -- These fields are currently not used
  signer_public_key VARCHAR(128),
  signer_weight INTEGER,
  signer_key VARCHAR(128),
  signer_type VARCHAR(128),

  --- This field is extracted from the corresponding operation
  created_at DATETIME
);
