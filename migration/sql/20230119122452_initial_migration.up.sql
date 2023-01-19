create table filemap (
   id bigint not null primary key auto_increment,
   filename varchar(255) not null, 
   filesize bigint not null, 
   collection varchar(255) not null, 
   created_at bigint not null
);

create table archivemap (
   id bigint not null primary key  auto_increment,
   filename varchar(255) not null, 
   filesize bigint not null,
   collection varchar(255) not null, 
   created_at bigint not null
);