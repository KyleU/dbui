select
  c.oid as "oid",
  n.nspname as "schema",
  c.relname as "name",
  case c.relkind when 'r' then 'table' when 'v' then 'view' when 'm' then 'materialized view' when 'i' then 'index' when 's' then 'sequence' when 's' then 'special' when 'f' then 'foreign table' when 'p' then 'table' when 'i' then 'index' end as "type",
  pg_catalog.pg_get_userbyid(c.relowner) as "owner",
  pg_catalog.pg_table_size(c.oid) as "size",
  pg_catalog.obj_description(c.oid, 'pg_class') as "description"
from pg_catalog.pg_class c
  left join pg_catalog.pg_namespace n on n.oid = c.relnamespace
where
  c.relkind in ('r','p','s','v','m','f','')
  and n.nspname !~ '^pg_toast'
  and n.nspname !~ '^pg_catalog'
  and pg_catalog.pg_table_is_visible(c.oid)
order by 2,3;

-- <%: func ListTables(buffer *bytes.Buffer) %>
