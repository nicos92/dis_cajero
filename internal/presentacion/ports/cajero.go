package ports

import (
	"github.com/nicos92/dis_cajero/internal/dominio/entidad"
)

type Cajero interface {
	ValidarTarjeta(tarjeta entidad.Tarjeta) error
	AutenticarUsuario(tarjeta entidad.Tarjeta, pin string) (entidad.Sesion, error)
	ObtenerCuentasAsociadas() ([]entidad.Cuenta, error)
	SeleccionarCuenta(numeroCuenta string) (entidad.Sesion, error)
	SesionActual() entidad.Sesion
	ConsultarSaldo() (float64, entidad.EstadoTransaccion, error)
	ExtraerEfectivo(monto float64) (entidad.EstadoTransaccion, int, error)
	RealizarTransferencia(cuentaDestino entidad.Cuenta, monto float64) (entidad.EstadoTransaccion, int, error)
	CerrarSesion() (bool, error)
}
