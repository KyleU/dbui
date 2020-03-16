create table if not exists projects (
  key varchar(512) not null primary key,
  title varchar(512) not null,
  description text,
  owner uuid,
  engine varchar(64) not null,
  url text not null
);

-- <%: func CreateTableProject(w io.Writer) %>
