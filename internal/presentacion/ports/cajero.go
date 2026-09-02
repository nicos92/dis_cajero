// internal/presentacion/ports.go (nuevo archivo)
package ports

import (
	"context"

	"github.com/nicos92/dis_cajero/internal/dominio/entidad"
)

type Cajero interface {
	AutenticarUsuario(ctx context.Context, tarjeta entidad.Tarjeta, cuenta entidad.Cuenta) (entidad.Sesion, error)
	ConsultarSaldo(ctx context.Context) (float64, entidad.EstadoTransaccion, error)
	ExtraerEfectivo(ctx context.Context, monto float64) (entidad.EstadoTransaccion, int, error)
	RealizarTransferencia(ctx context.Context, cuentaDestino entidad.Cuenta, monto float64) (entidad.EstadoTransaccion, int, error)
	CerrarSesion(ctx context.Context) (bool, error)
}
