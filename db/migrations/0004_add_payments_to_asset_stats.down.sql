ALTER TABLE asset_stats DROP COLUMN payments;
ALTER TABLE asset_stats ADD COLUMN num_effects integer;

CREATE OR REPLACE FUNCTION repopulate_asset_stats()
  RETURNS VOID
AS
$$
DECLARE
   t_row effects%rowtype;
BEGIN
    TRUNCATE asset_stats;
    FOR t_row in SELECT * FROM effects LOOP
      INSERT INTO asset_stats(paging_token, asset_code, asset_issuer, asset_type, created_at, total_amount, num_accounts, num_effects)
          VALUES (t_row.paging_token, t_row.asset_code, t_row.asset_issuer, t_row.asset_type, t_row.created_at,
              (SELECT COALESCE(SUM(amount), 0) FROM effects WHERE type='account_debited' AND account=t_row.asset_issuer AND effect_id <= t_row.effect_id),
              (SELECT COUNT(DISTINCT account) FROM effects WHERE effect_id <= t_row.effect_id),
              (SELECT COUNT(*) FROM effects WHERE effect_id <= t_row.effect_id)
          );
    END LOOP;
END;
$$
LANGUAGE plpgsql;

-- Re-populate asset_stats table
SELECT repopulate_asset_stats();

DROP FUNCTION repopulate_asset_stats;
