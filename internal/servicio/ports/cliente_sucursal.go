package ports

import (
	entidad "github.com/nicos92/dis_cajero/internal/dominio/entidad"
)

type ClienteSucursal interface {
	ValidarTajeta(tarjeta entidad.Tarjeta) (bool, error)
}
