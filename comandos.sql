CREATE DATABASE proyecto
USE proyecto
CREATE TABLE medicamento(
id SERIAL PRIMARY KEY,
nombre VARCHAR(200),
principio_activo VARCHAR(200),
presentacion VARCHAR(200),
precio FLOAT 

)
SELECT * FROM medicamento