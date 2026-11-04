package ports

import entidad "github.com/nicos92/dis_cajero/internal/dominio/entidad"

type ClienteRedCentral interface {
	ObtenerCuenta(numero string) (entidad.Cuenta, error)
	ObtenerCuentas(numeros []string) ([]entidad.Cuenta, error)
	ConsultarSaldo(cuenta entidad.Cuenta) (float64, error)
	DebitarExtraccion(cuenta entidad.Cuenta, monto float64) (entidad.EstadoTransaccion, int, error)
	TransferirACuenta(cuentaOrigen entidad.Cuenta, cuentaDestino entidad.Cuenta, monto float64) (entidad.EstadoTransaccion, int, error)
}
