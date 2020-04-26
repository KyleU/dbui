drop table if exists pktest;

create table pktest (
  id integer,
  foo varchar(100),
  bar uuid unique,
  primary key (id, foo)
);

insert into pktest (
  id, foo, bar
) values (
  0, 'a', '00000000-0000-0000-0000-000000000000'
);

-- <%: func ExamplePKTest(w io.Writer) %>
