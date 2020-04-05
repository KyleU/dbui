select
  m.name as table_name,
  p.name as column_name,
  p.type as data_type,
  p."notnull" as is_nullable,
  p.dflt_value as column_default
from
  sqlite_master as m join pragma_table_info(m.name) as p
order by
  m.name,
  p.cid;

-- <%: func ListColumnsSQLite(w io.Writer) %>
