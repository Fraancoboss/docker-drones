// Archivo: tools/drone-observe/main.go
// Rol: punto de entrada del CLI drone-observe para validar el pipeline de observabilidad.
// No hace: logica de negocio ni scraping directo fuera de los comandos.
package main

import (
	"os"

	"drone-observe/cmd"
)

func main() {
	os.Exit(cmd.Execute(os.Args))
}
