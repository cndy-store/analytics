-- NOTE: Taken from migration 0003
DROP TABLE asset_stats;
CREATE TABLE asset_stats (
  paging_token character varying(64) PRIMARY KEY,
  asset_type character varying(64),
  asset_code character varying(12),
  asset_issuer character varying(56),
  total_amount bigint,
  num_accounts integer,
  payments integer,
  created_at timestamp without time zone
);

-- NOTE: Taken from migration 0004
CREATE OR REPLACE FUNCTION repopulate_asset_stats()
  RETURNS VOID
AS
$$
DECLARE
   t_row effects%rowtype;
BEGIN
    TRUNCATE asset_stats;
    FOR t_row in SELECT * FROM effects LOOP
      INSERT INTO asset_stats(paging_token, asset_code, asset_issuer, asset_type, created_at, total_amount, num_accounts, payments)
      VALUES (t_row.paging_token, t_row.asset_code, t_row.asset_issuer, t_row.asset_type, t_row.created_at,
          (SELECT COALESCE(SUM(amount), 0) FROM effects WHERE type='account_debited' AND account=t_row.asset_issuer AND effect_id <= t_row.effect_id),
          (SELECT COUNT(DISTINCT account) FROM effects WHERE account != t_row.asset_issuer AND effect_id <= t_row.effect_id),
          (SELECT COUNT(*) FROM effects WHERE type='account_debited' AND account != t_row.asset_issuer AND effect_id <= t_row.effect_id)
      );
    END LOOP;
END;
$$
LANGUAGE plpgsql;

-- Re-populate asset_stats table
SELECT repopulate_asset_stats();
