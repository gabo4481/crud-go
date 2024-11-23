use proyecto;
create table medicamento(
id serial primary key,
nombre varchar(200) not null,
principio_activo varchar(200) not null,
presentacion varchar(200) not null,
precio float not null
)

select * from medicamento