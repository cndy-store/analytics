CREATE OR REPLACE FUNCTION asset_stats(asset_code_filter character varying(12), asset_issuer_filter character varying(56))
  RETURNS TABLE (
    -- TODO: ID to order required?
    effect_id character varying(56),
    paging_token character varying(64),
    asset_type character varying(64),
    asset_code character varying(12),
    asset_issuer character varying(56),
    issued bigint,
    transferred bigint,
    accounts_with_trustline integer,
    accounts_with_payments integer,
    payments integer,
    created_at timestamp without time zone
  )
AS
$$
#variable_conflict use_column
DECLARE
   t_row record;
BEGIN
  FOR t_row in SELECT effect_id, paging_token, asset_type, asset_code, asset_issuer, created_at FROM effects ORDER BY created_at LOOP
    -- Next if asset_code and asset_issuer do not match
    CONTINUE WHEN t_row.asset_code != asset_code_filter AND t_row.asset_issuer != asset_issuer_filter;

    effect_id               := t_row.effect_id;
    paging_token            := t_row.paging_token;
    asset_type              := t_row.asset_type;
    asset_code              := t_row.asset_code;
    asset_issuer            := t_row.asset_issuer;
    created_at              := t_row.created_at;
    issued                  := (SELECT COALESCE(SUM(amount), 0) FROM effects WHERE type = 'account_debited' AND account = t_row.asset_issuer AND effect_id <= t_row.effect_id);
    transferred             := (SELECT COALESCE(SUM(amount), 0) FROM effects WHERE type = 'account_debited' AND account != t_row.asset_issuer AND effect_id <= t_row.effect_id);
    accounts_with_trustline := (SELECT COUNT(DISTINCT account) FROM effects WHERE account != t_row.asset_issuer AND type = 'trustline_created' AND effect_id <= t_row.effect_id);
    accounts_with_payments  := (SELECT COUNT(DISTINCT account) FROM effects WHERE account != t_row.asset_issuer AND type = 'account_debited' AND effect_id <= t_row.effect_id);
    payments                := (SELECT COUNT(*) FROM effects WHERE type = 'account_debited' AND account != t_row.asset_issuer AND effect_id <= t_row.effect_id);
    RETURN NEXT;
  END LOOP;
END;
$$
LANGUAGE plpgsql;

DROP TABLE asset_stats;
DROP FUNCTION repopulate_asset_stats;
