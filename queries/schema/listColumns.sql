select
  c.table_schema,
  c.table_name,
  c.column_name,
  c.ordinal_position,
  c.column_default,
  c.is_nullable,
  c.data_type,
  e.data_type as array_type,
  c.character_maximum_length,
  c.character_octet_length,
  c.numeric_precision,
  c.numeric_precision_radix,
  c.numeric_scale,
  c.datetime_precision,
  c.interval_type,
  c.domain_schema,
  c.domain_name,
  c.udt_schema,
  c.udt_name,
  c.dtd_identifier,
  c.is_updatable
from information_schema.columns c left join information_schema.element_types e on (
    (c.table_catalog, c.table_schema, c.table_name, 'TABLE', c.dtd_identifier) = (e.object_catalog, e.object_schema, e.object_name, e.object_type, e.collection_type_identifier)
  )
where table_schema = 'public'
order by c.table_schema, c.table_name, c.ordinal_position;

-- <%: func ListColumns(w io.Writer) %>
