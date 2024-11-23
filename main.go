package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "filippo.io/edwards25519"
	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "proyecto"

	conexion, err := sql.Open(Driver, Usuario+":"+"@tcp(127.0.0.1)/"+Contrasenia+Nombre)
	if err != nil {
		panic(err.Error())
	}

	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

type Medicamentos struct {
	Id               int64
	Nombre           string
	Principio_activo string
	Presentacion     string
	Precio           float64
}

func main() {
	http.Handle("/imagenes/", http.StripPrefix("/imagenes/", http.FileServer(http.Dir("./imagenes"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)
	http.HandleFunc("/cancelar", Cancelar)
	http.HandleFunc("/version", Version)

	fmt.Println("Servidor Corriendo...")
	fmt.Println("run Server: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func Inicio(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w,"hola mundo")

	conexion_establecida := conexionBD()
	registros, err := conexion_establecida.Query("SELECT * FROM medicamento")
	if err != nil {
		panic(err.Error())
	}
	medicamento := Medicamentos{}
	arreglo_medicamentos := []Medicamentos{}

	for registros.Next() {
		var id int
		var nombre, principio_activo, presentacion string
		var precio float64
		err := registros.Scan(&id, &nombre, &principio_activo, &presentacion, &precio)
		if err != nil {
			panic(err.Error())
		}
		medicamento.Id = int64(id)
		medicamento.Nombre = nombre
		medicamento.Principio_activo = principio_activo
		medicamento.Presentacion = presentacion
		medicamento.Precio = precio

		arreglo_medicamentos = append(arreglo_medicamentos, medicamento)
	}

	plantillas.ExecuteTemplate(w, "inicio", arreglo_medicamentos)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		principio_activo := r.FormValue("principio")
		presentacion := r.FormValue("presentacion")
		precio := r.FormValue("precio")

		conexion_establecida := conexionBD()
		insertar_registros, err := conexion_establecida.Prepare("INSERT INTO medicamento(nombre,principio_activo,presentacion,precio) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insertar_registros.Exec(nombre, principio_activo, presentacion, precio)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	id_medicamento := r.URL.Query().Get("id")

	conexion_establecida := conexionBD()
	insertar_registros, err := conexion_establecida.Prepare("DELETE FROM medicamento where id=?")
	if err != nil {
		panic(err.Error())
	}
	insertar_registros.Exec(id_medicamento)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func Editar(w http.ResponseWriter, r *http.Request) {
	id_medicamento := r.URL.Query().Get("id")

	conexion_establecida := conexionBD()
	registros_editar, err := conexion_establecida.Query("SELECT * FROM medicamento where id=?", id_medicamento)
	if err != nil {
		panic(err.Error())
	}
	medicamento := Medicamentos{}

	for registros_editar.Next() {
		var id int
		var nombre, principio_activo, presentacion string
		var precio float64
		err := registros_editar.Scan(&id, &nombre, &principio_activo, &presentacion, &precio)
		if err != nil {
			panic(err.Error())
		}
		medicamento.Id = int64(id)
		medicamento.Nombre = nombre
		medicamento.Principio_activo = principio_activo
		medicamento.Presentacion = presentacion
		medicamento.Precio = precio
	}

	plantillas.ExecuteTemplate(w, "editar", medicamento)
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		principio_activo := r.FormValue("principio")
		presentacion := r.FormValue("presentacion")
		precio := r.FormValue("precio")

		conexion_establecida := conexionBD()
		actualizar_registros, err := conexion_establecida.Prepare("UPDATE medicamento SET nombre=?,principio_activo=?,presentacion=?,precio=? where id=?")
		if err != nil {
			panic(err.Error())
		}
		actualizar_registros.Exec(nombre, principio_activo, presentacion, precio, id)
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Cancelar(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Version(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "version", nil)
}
