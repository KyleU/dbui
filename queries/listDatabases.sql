select
  d.datname as "name",
  pg_catalog.pg_get_userbyid(d.datdba) as "owner",
  pg_catalog.pg_encoding_to_char(d.encoding) as "encoding",
  d.datcollate as "collate",
  d.datctype as "ctype",
  pg_catalog.array_to_string(d.datacl, E'\n') AS "access"
from pg_catalog.pg_database d
order by 1;

-- <%: func ListDatabases(buffer *bytes.Buffer) %>
