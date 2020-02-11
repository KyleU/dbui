select i.relname as idx,
  a.attname as tn,
  idx.indisprimary as pk,
  idx.indisunique as u,
  array_position(idx.indkey, a.attnum) as seq
from pg_class t,
  pg_class i,
  pg_index idx,
  pg_attribute a,
  pg_namespace n
where
  t.oid = idx.indrelid
  and i.oid = idx.indexrelid
  and n.oid = t.relnamespace
  and a.attrelid = t.oid
  and a.attnum = any(idx.indkey)
  and t.relkind = 'r'
  and n.nspname = current_schema()
order by idx, seq;

-- <%: func ListIndexes(w io.Writer) %>
