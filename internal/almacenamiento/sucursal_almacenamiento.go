package almacenamiento

import (
	"fmt"
	"time"

	entidad "github.com/nicos92/dis_cajero/internal/dominio/entidad"
)

type SucursalAlmacenamiento struct {
	Tarjetas map[string]entidad.Tarjeta
}

func (sucursal *SucursalAlmacenamiento) ValidarTarjeta(tarjeta entidad.Tarjeta) (entidad.Tarjeta, error) {
	tarjetaEncontrada, existe := sucursal.Tarjetas[tarjeta.Numero]

	if !existe {
		return entidad.Tarjeta{}, fmt.Errorf("%w: %s", entidad.ErrTarjetaInvalida, tarjeta.Numero)
	}

	if tarjetaEncontrada.EstaVencida(time.Now()) {
		return entidad.Tarjeta{}, fmt.Errorf("%w: %s", entidad.ErrTarjetaVencida, tarjeta.Numero)
	}

	return tarjetaEncontrada, nil
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
		Numero:           "123456",
		Pin:              "1111",
		Vencimiento:      time.Now().AddDate(1, 0, 0),
		CuentasAsociadas: []string{"1111", "333"},
	}
	tarjetaVencida := entidad.Tarjeta{
		Numero:           "999999",
		Pin:              "2222",
		Vencimiento:      time.Now().AddDate(0, 0, -1),
		CuentasAsociadas: []string{"2222"},
	}

	sucursal.Tarjetas[tarjetaValida.Numero] = tarjetaValida
	sucursal.Tarjetas[tarjetaVencida.Numero] = tarjetaVencida
}
