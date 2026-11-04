package almacenamiento

import (
	"fmt"
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
		Saldo:  234.234,
		Activa: true,
	}
	cuentaDos := entidad.Cuenta{
		Numero: "2222",
		Saldo:  543.54,
		Activa: false,
	}
	cuentaTres := entidad.Cuenta{
		Numero: "333",
		Saldo:  543.54,
		Activa: true,
	}

	redCentral.Cuentas[cuentaUno.Numero] = cuentaUno
	redCentral.Cuentas[cuentaDos.Numero] = cuentaDos
	redCentral.Cuentas[cuentaTres.Numero] = cuentaTres
}

func (redCentral *RedCentralAlmacenamiento) ObtenerCuenta(numero string) (entidad.Cuenta, error) {
	cuentaGuardada, existe := redCentral.Cuentas[numero]

	if !existe {
		return entidad.Cuenta{}, fmt.Errorf("%w: %s", entidad.ErrCuentaNoExiste, numero)
	}

	return cuentaGuardada, nil
}

func (redCentral *RedCentralAlmacenamiento) ObtenerCuentas(numeros []string) ([]entidad.Cuenta, error) {
	cuentas := make([]entidad.Cuenta, 0, len(numeros))

	for _, numero := range numeros {
		cuenta, err := redCentral.ObtenerCuenta(numero)
		if err != nil {
			return nil, err
		}

		cuentas = append(cuentas, cuenta)
	}

	return cuentas, nil
}

func (redCentral *RedCentralAlmacenamiento) ConsultarSaldo(cuenta entidad.Cuenta) (float64, error) {
	cuentaGuardada, err := redCentral.ObtenerCuenta(cuenta.Numero)
	if err != nil {
		return 0, err
	}

	return cuentaGuardada.Saldo, nil
}

func (redCentral *RedCentralAlmacenamiento) DebitarExtraccion(cuenta entidad.Cuenta, monto float64) (entidad.EstadoTransaccion, int, error) {
	cuentaGuardada, err := redCentral.ObtenerCuenta(cuenta.Numero)
	if err != nil {
		return entidad.EstadoCancelado, 0, err
	}

	if !cuentaGuardada.Activa {
		return entidad.EstadoCancelado, 0, fmt.Errorf("%w: %s", entidad.ErrCuentaInactiva, cuenta.Numero)
	}

	if cuentaGuardada.Saldo < monto {
		return entidad.EstadoCancelado, 0, entidad.ErrSaldoInsuficiente
	}

	cuentaGuardada.Saldo -= monto

	redCentral.Cuentas[cuentaGuardada.Numero] = cuentaGuardada

	id := len(redCentral.Transacciones) + 1
	redCentral.Transacciones[strconv.Itoa(id)] = id

	return entidad.EstadoConfirmado, id, nil
}

func (redCentral *RedCentralAlmacenamiento) TransferirACuenta(origen entidad.Cuenta, destino entidad.Cuenta, monto float64) (entidad.EstadoTransaccion, int, error) {
	cuentaOrigen, err := redCentral.ObtenerCuenta(origen.Numero)
	if err != nil {
		return entidad.EstadoCancelado, 0, err
	}

	if !cuentaOrigen.Activa {
		return entidad.EstadoCancelado, 0, fmt.Errorf("%w: %s", entidad.ErrCuentaInactiva, origen.Numero)
	}

	if cuentaOrigen.Saldo < monto {
		return entidad.EstadoCancelado, 0, entidad.ErrSaldoInsuficiente
	}

	cuentaDestino, err := redCentral.ObtenerCuenta(destino.Numero)
	if err != nil {
		return entidad.EstadoCancelado, 0, err
	}

	if !cuentaDestino.Activa {
		return entidad.EstadoCancelado, 0, fmt.Errorf("%w: %s", entidad.ErrCuentaInactiva, destino.Numero)
	}

	if cuentaDestino.Numero == cuentaOrigen.Numero {
		return entidad.EstadoCancelado, 0, fmt.Errorf("%w: no puede transferir a la misma cuenta", entidad.ErrCuentaDestinoInvalida)
	}

	cuentaOrigen.Saldo -= monto
	cuentaDestino.Saldo += monto

	redCentral.Cuentas[cuentaOrigen.Numero] = cuentaOrigen
	redCentral.Cuentas[cuentaDestino.Numero] = cuentaDestino

	id := len(redCentral.Transacciones) + 1
	redCentral.Transacciones[strconv.Itoa(id)] = id

	return entidad.EstadoConfirmado, id, nil
}
