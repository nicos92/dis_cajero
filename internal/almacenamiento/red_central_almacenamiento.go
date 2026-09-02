package almacenamiento

import (
	"errors"
	"strconv"

	entidad "github.com/nicos92/dis_cajero/internal/dominio/entidad"
)

type RedCentralAlmacenamiento struct {
	Cuentas       map[string]entidad.Cuenta
	Transacciones map[string]int
}

func NuevoRedCentralAlmacenamiento() *RedCentralAlmacenamiento {
	redCentral := &RedCentralAlmacenamiento{
		Cuentas:       make(map[string]entidad.Cuenta),
		Transacciones: make(map[string]int),
	}

	redCentral.seedCuentas()

	return redCentral
}

func (redCentral *RedCentralAlmacenamiento) seedCuentas() {
	cuentaUno := entidad.Cuenta{
		Numero: "1111",
		Pin:    "1111",
		Saldo:  234.234,
		Activa: true,
	}
	cuentaDos := entidad.Cuenta{
		Numero: "2222",
		Pin:    "2222",
		Saldo:  543.54,
		Activa: false,
	}
	cuentaTres := entidad.Cuenta{
		Numero: "333",
		Pin:    "2223332",
		Saldo:  543.54,
		Activa: true,
	}

	redCentral.Cuentas[cuentaUno.Numero] = cuentaUno
	redCentral.Cuentas[cuentaDos.Numero] = cuentaDos
	redCentral.Cuentas[cuentaTres.Numero] = cuentaTres

}

func (redCentral *RedCentralAlmacenamiento) ValidarPin(cuenta entidad.Cuenta) (bool, error) {
	cuentaGuardada, existe := redCentral.Cuentas[cuenta.Numero]

	if !existe {
		return false, errors.New("La cuenta no existe")
	}

	if !cuentaGuardada.Activa {
		return false, errors.New("La cuenta no esta activada")
	}

	if cuentaGuardada.Pin != cuenta.Pin {
		return false, errors.New("El Pin es incorrecto")
	}

	return true, nil
}

func (redCentral *RedCentralAlmacenamiento) ConsultarSaldo(cuenta entidad.Cuenta) (float64, error) {
	cuentaGuardada, existe := redCentral.Cuentas[cuenta.Numero]

	if !existe {
		return 0, errors.New("La cuenta no existe")
	}

	return cuentaGuardada.Saldo, nil
}

func (redCentral *RedCentralAlmacenamiento) DebitarExtraccion(cuenta entidad.Cuenta, monto float64) (entidad.EstadoTransaccion, int, error) {
	cuentaGuardada, existe := redCentral.Cuentas[cuenta.Numero]

	if !existe {
		return entidad.EstadoCancelado, 0, errors.New("La cuenta no existe")
	}

	if cuentaGuardada.Saldo < monto {
		return entidad.EstadoCancelado, 0, errors.New("Saldo insuficiente para extraer")
	}

	cuentaGuardada.Saldo -= monto

	redCentral.Cuentas[cuentaGuardada.Numero] = cuentaGuardada

	id := len(redCentral.Transacciones) + 1
	redCentral.Transacciones[strconv.Itoa(id)] = id

	return entidad.EstadoConfirmado, id, nil
}

func (redCentral *RedCentralAlmacenamiento) TransferirACuenta(origen entidad.Cuenta, destino entidad.Cuenta, monto float64) (entidad.EstadoTransaccion, int, error) {
	cuentaOrigen, existe := redCentral.Cuentas[origen.Numero]

	if !existe {
		return entidad.EstadoCancelado, 0, errors.New("La cuenta de origen no existe")
	}

	if cuentaOrigen.Saldo < monto {
		return entidad.EstadoCancelado, 0, errors.New("Saldo insuficiente para transferir")
	}

	cuentaDestino, existe := redCentral.Cuentas[destino.Numero]

	if !existe {
		return entidad.EstadoCancelado, 0, errors.New("La cuenta de destino no existe")
	}

	if !cuentaDestino.Activa {
		return entidad.EstadoCancelado, 0, errors.New("La cuenta de destino no está activada.")

	}

	if cuentaDestino.Numero == cuentaOrigen.Numero {
		return entidad.EstadoCancelado, 0, errors.New("La cuenta de destino no puede ser su misma cuenta")

	}
	cuentaOrigen.Saldo -= monto
	cuentaDestino.Saldo += monto

	redCentral.Cuentas[cuentaOrigen.Numero] = cuentaOrigen

	redCentral.Cuentas[cuentaDestino.Numero] = cuentaDestino

	id := len(redCentral.Transacciones) + 1
	redCentral.Transacciones[strconv.Itoa(id)] = id

	return entidad.EstadoConfirmado, id, nil
}
