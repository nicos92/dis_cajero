package entidad

import "time"

type Tarjeta struct {
	Numero           string
	Pin              string
	Vencimiento      time.Time
	CuentasAsociadas []string
}

func (tarjeta Tarjeta) EstaVencida(ahora time.Time) bool {
	if tarjeta.Vencimiento.IsZero() {
		return false
	}

	return tarjeta.Vencimiento.Before(ahora)
}

func (tarjeta Tarjeta) TieneCuenta(numeroCuenta string) bool {
	for _, numero := range tarjeta.CuentasAsociadas {
		if numero == numeroCuenta {
			return true
		}
	}

	return false
}
