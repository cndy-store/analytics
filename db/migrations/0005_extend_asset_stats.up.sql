ALTER TABLE asset_stats DROP COLUMN total_amount;
ALTER TABLE asset_stats DROP COLUMN num_accounts;
ALTER TABLE asset_stats ADD COLUMN issued bigint;
ALTER TABLE asset_stats ADD COLUMN transferred bigint;
ALTER TABLE asset_stats ADD COLUMN accounts_with_trustline integer;
ALTER TABLE asset_stats ADD COLUMN accounts_with_payments integer;

CREATE OR REPLACE FUNCTION repopulate_asset_stats()
  RETURNS VOID
AS
$$
DECLARE
   t_row record;
BEGIN
    TRUNCATE asset_stats;
    FOR t_row in SELECT paging_token, effect_id, asset_type, asset_code, asset_issuer, created_at FROM effects ORDER BY effect_id LOOP
        INSERT INTO asset_stats(paging_token, asset_code, asset_issuer, asset_type, created_at, issued, transferred, accounts_with_trustline, accounts_with_payments, payments)
        VALUES (t_row.paging_token, t_row.asset_code, t_row.asset_issuer, t_row.asset_type, t_row.created_at,
            (SELECT COALESCE(SUM(amount), 0) FROM effects WHERE type = 'account_debited' AND account = t_row.asset_issuer AND effect_id <= t_row.effect_id),
            (SELECT COALESCE(SUM(amount), 0) FROM effects WHERE type = 'account_debited' AND account != t_row.asset_issuer AND effect_id <= t_row.effect_id),
            (SELECT COUNT(DISTINCT account) FROM effects WHERE account != t_row.asset_issuer AND type = 'trustline_created' AND effect_id <= t_row.effect_id),
            (SELECT COUNT(DISTINCT account) FROM effects WHERE account != t_row.asset_issuer AND type = 'account_debited' AND effect_id <= t_row.effect_id),
            (SELECT COUNT(*) FROM effects WHERE type = 'account_debited' AND account != t_row.asset_issuer AND effect_id <= t_row.effect_id)
        );
    END LOOP;
END;
$$
LANGUAGE plpgsql;

-- Re-populate asset_stats table
SELECT repopulate_asset_stats();
