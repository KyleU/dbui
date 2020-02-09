package query

/// List all databases
const ListDatabases = `select
  d.datname as "name",
  pg_catalog.pg_get_userbyid(d.datdba) as "owner",
  pg_catalog.pg_encoding_to_char(d.encoding) as "encoding",
  d.datcollate as "collate",
  d.datctype as "ctype",
  pg_catalog.array_to_string(d.datacl, E'\n') AS "access"
from pg_catalog.pg_database d
order by 1;`

/// List all tables within a database
const ListTables = `select
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
order by 2,3;`

/// List all columns within a database's table
const ListColumns = `select
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
order by c.table_schema, c.table_name, c.ordinal_position;`
