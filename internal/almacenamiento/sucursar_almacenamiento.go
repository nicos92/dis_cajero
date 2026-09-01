package almacenamiento

import (
	"errors"
	"time"

	entidad "github.com/nicos92/dis_cajero/internal/dominio/entidad"
)

type SucursalAlmacenamiento struct {
	Tarjetas map[string]entidad.Tarjeta
}

func (sucursal *SucursalAlmacenamiento) ValidarTarjeta(tajeta entidad.Tarjeta) (bool, error) {

	_, existe := sucursal.Tarjetas[tajeta.Numero]

	if !existe {
		return false, errors.New("Tarjeta invalida")
	}

	if tajeta.Vencimiento.Before(time.Now()) {
		return false, errors.New("Tarjeta vencida.")
	}

	return true, nil
}
