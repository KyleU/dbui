select
  t.relname as table_name,
  i.relname as index_name,
  idx.indisprimary as pk,
  idx.indisunique as u,
  array_to_string(array_agg(a.attname), ',') as column_names
from
  pg_class t,
  pg_class i,
  pg_index idx,
  pg_attribute a,
  pg_namespace n
where
    t.oid = idx.indrelid
    and i.oid = idx.indexrelid
    and a.attrelid = t.oid
    and n.oid = t.relnamespace
    and a.attnum = any(idx.indkey)
    and t.relkind = 'r'
    and n.nspname = current_schema()
group by
  t.relname,
  i.relname,
  idx.indisprimary,
  idx.indisunique
order by
  t.relname,
  i.relname;

-- <%: func ListIndexesPostgres(w io.Writer) %>
