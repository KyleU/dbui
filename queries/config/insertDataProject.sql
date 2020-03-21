insert into projects (
  key,
  title,
  description,
  owner,
  engine,
  url
) values (
  'config',
  'Config',
  'Internal config database',
  '00000000-0000-0000-0000-000000000000',
  'sqlite',
  'dbui.db'
);

-- <%: func InsertDataProject(w io.Writer) %>
