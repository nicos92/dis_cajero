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

func NuevoSucursalAlmacenamiento() *SucursalAlmacenamiento {
	sucursal := &SucursalAlmacenamiento{
		Tarjetas: make(map[string]entidad.Tarjeta),
	}

	sucursal.seedTarjetas()

	return sucursal
}

func (sucursal *SucursalAlmacenamiento) seedTarjetas() {

	tarjetaValida := entidad.Tarjeta{
		Numero:      "123456",
		Vencimiento: time.Now().AddDate(1, 0, 0),
	}
	tarjetaVencida := entidad.Tarjeta{
		Numero:      "999999",
		Vencimiento: time.Now().AddDate(0, 0, -1),
	}

	sucursal.Tarjetas[tarjetaValida.Numero] = tarjetaValida
	sucursal.Tarjetas[tarjetaVencida.Numero] = tarjetaVencida
}
