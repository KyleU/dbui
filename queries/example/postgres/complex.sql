
create extension if not exists hstore;

drop table if exists complex;

create table complex (
  id serial primary key,
  ay bytea,
  ab bool[],
  ao oid[],
  af float8[],
  aj json[],
  at text[],
  ats timestamp[],
  h hstore
);

insert into complex (
  ay, ab, ao, af, aj, at, ats, h
) values (
  '{ 0, 1, 2 }', '{true, false}', '{0, 1}', '{0.1, 0.1}', '{ "{}", "{}" }', '{ "a", "b", "c" }', '{ 2000-01-01, 2000-01-02 }', 'a => x, b => y, c => z'
);

-- <%: func ExampleComplex(w io.Writer) %>
