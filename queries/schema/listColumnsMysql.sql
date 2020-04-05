select
  table_schema as table_schema,
  table_name as table_name,
  column_name as column_name,
  ordinal_position as ordinal_position,
  column_default as column_default,
  is_nullable as is_nullable,
  data_type as data_type,
  numeric_precision as numeric_precision,
  numeric_scale as numeric_scale,
  character_maximum_length as character_maximum_length,
  datetime_precision as datetime_precision
from
  information_schema.columns
where
  table_schema = 'dbui'
order by
  table_schema,
  table_name,
  ordinal_position;

-- <%: func ListColumnsMySQL(w io.Writer) %>
