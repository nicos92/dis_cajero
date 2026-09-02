package main

import (
	"github.com/nicos92/dis_cajero/internal/almacenamiento"
	"github.com/nicos92/dis_cajero/internal/dominio/entidad"
	"github.com/nicos92/dis_cajero/internal/presentacion"
	"github.com/nicos92/dis_cajero/internal/servicio"
)

func main() {
	sucursal := almacenamiento.NuevoSucursalAlmacenamiento()
	redCentral := almacenamiento.NuevoRedCentralAlmacenamiento()

	cajero := &servicio.CajeroService{
		ClienteSucursal:   sucursal,
		ClienteRedCentral: redCentral,
		Sesion:            entidad.Sesion{},
	}

	consola := presentacion.NuevaConsola(cajero)

	consola.Iniciar()

}
