CREATE TABLE assets (
  -- Code and issuer combination needs to be unique, can serve as primary key
  code character varying(12) not null,
  issuer character varying(56) not null,
  primary key (code, issuer),

  type character varying(64) default 'credit_alphanum4',
  created_at timestamp without time zone default current_timestamp
);

-- Register CNDY
INSERT INTO assets(type, code, issuer) VALUES(
  'credit_alphanum4',
  'CNDY',
  'GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX'
);
