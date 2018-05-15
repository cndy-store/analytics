CREATE TABLE cursors (
  id SERIAL PRIMARY KEY,
  paging_token character varying(64) NOT NULL
);

-- Set genesis cursor (Mon Mar 12 18:01:12 CET 2018)
INSERT INTO cursors(paging_token) VALUES('33819440072110101-2');
