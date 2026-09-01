package ports

import "github.com/nicos92/dis_cajero/internal/domain/entity"

type ClienteRedCentral interface {
	ValidarPin(cuenta entity.Cuenta) (bool, error)
	ConsultarSaldo(cuenta entity.Cuenta) (float64, error)
	DebitarExtaccion(cuenta entity.Cuenta, monto float64) (entity.EstadoTransaccion, int, error)
	TransferirACuenta(cuentaOrigen entity.Cuenta, cuentaDestino entity.Cuenta, monto float64) (entity.EstadoTransaccion, int, error)
}
