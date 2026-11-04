package servicio

import (
	"fmt"

	entidad "github.com/nicos92/dis_cajero/internal/dominio/entidad"
	"github.com/nicos92/dis_cajero/internal/servicio/ports"
)

type CajeroService struct {
	ClienteRedCentral ports.ClienteRedCentral
	ClienteSucursal   ports.ClienteSucursal
	Sesion            entidad.Sesion
}

func (cajeroService *CajeroService) ValidarTarjeta(tarjeta entidad.Tarjeta) error {
	_, err := cajeroService.ClienteSucursal.ValidarTarjeta(tarjeta)
	return err
}

func (cajeroService *CajeroService) AutenticarUsuario(tarjeta entidad.Tarjeta, pin string) (entidad.Sesion, error) {
	tarjetaRegistrada, err := cajeroService.ClienteSucursal.ValidarTarjeta(tarjeta)
	if err != nil {
		return entidad.Sesion{}, err
	}

	if pin == "" {
		return entidad.Sesion{}, fmt.Errorf("%w: el pin es obligatorio", entidad.ErrCredencialesInvalidas)
	}

	if tarjetaRegistrada.Pin != pin {
		return entidad.Sesion{}, fmt.Errorf("%w: pin incorrecto", entidad.ErrPinIncorrecto)
	}

	cajeroService.Sesion = entidad.Sesion{
		Tarjeta: tarjetaRegistrada,
		Activa:  false,
	}

	return cajeroService.Sesion, nil
}

func (cajeroService *CajeroService) ObtenerCuentasAsociadas() ([]entidad.Cuenta, error) {
	if cajeroService.Sesion.Tarjeta.Numero == "" {
		return nil, entidad.ErrSesionInactiva
	}

	cuentas, err := cajeroService.ClienteRedCentral.ObtenerCuentas(cajeroService.Sesion.Tarjeta.CuentasAsociadas)
	if err != nil {
		return nil, err
	}

	if len(cuentas) == 0 {
		return nil, entidad.ErrCuentasNoAsociadas
	}

	return cuentas, nil
}

func (cajeroService *CajeroService) SeleccionarCuenta(numeroCuenta string) (entidad.Sesion, error) {
	if cajeroService.Sesion.Tarjeta.Numero == "" {
		return entidad.Sesion{}, entidad.ErrSesionInactiva
	}

	if !cajeroService.Sesion.Tarjeta.TieneCuenta(numeroCuenta) {
		return entidad.Sesion{}, fmt.Errorf("%w: la cuenta %s no está asociada a la tarjeta", entidad.ErrCuentaNoExiste, numeroCuenta)
	}

	cuenta, err := cajeroService.ClienteRedCentral.ObtenerCuenta(numeroCuenta)
	if err != nil {
		return entidad.Sesion{}, err
	}

	if !cuenta.Activa {
		return entidad.Sesion{}, fmt.Errorf("%w: %s", entidad.ErrCuentaInactiva, numeroCuenta)
	}

	cajeroService.Sesion.CuentaCliente = cuenta
	cajeroService.Sesion.Activa = true

	return cajeroService.Sesion, nil
}

func (cajeroService *CajeroService) SesionActual() entidad.Sesion {
	return cajeroService.Sesion
}

func (cajeroService *CajeroService) ConsultarSaldo() (float64, entidad.EstadoTransaccion, error) {
	if !cajeroService.Sesion.Activa {
		return 0, entidad.EstadoCancelado, entidad.ErrSesionInactiva
	}

	saldo, err := cajeroService.ClienteRedCentral.ConsultarSaldo(cajeroService.Sesion.CuentaCliente)
	if err != nil {
		return 0, entidad.EstadoCancelado, err
	}

	return saldo, entidad.EstadoConfirmado, nil
}

func (cajeroService *CajeroService) ExtraerEfectivo(monto float64) (entidad.EstadoTransaccion, int, error) {
	if !cajeroService.Sesion.Activa {
		return entidad.EstadoCancelado, 0, entidad.ErrSesionInactiva
	}

	estado, id, err := cajeroService.ClienteRedCentral.DebitarExtraccion(cajeroService.Sesion.CuentaCliente, monto)
	if err != nil {
		return estado, id, err
	}

	return estado, id, nil
}

func (cajeroService *CajeroService) RealizarTransferencia(cuentaDestino entidad.Cuenta, monto float64) (entidad.EstadoTransaccion, int, error) {
	if !cajeroService.Sesion.Activa {
		return entidad.EstadoCancelado, 0, entidad.ErrSesionInactiva
	}

	estado, id, err := cajeroService.ClienteRedCentral.TransferirACuenta(cajeroService.Sesion.CuentaCliente, cuentaDestino, monto)
	if err != nil {
		return estado, id, err
	}

	return estado, id, nil
}

func (cajeroService *CajeroService) CerrarSesion() (bool, error) {
	if !cajeroService.Sesion.Activa {
		return false, entidad.ErrSesionInactiva
	}

	cajeroService.Sesion = entidad.Sesion{}

	return cajeroService.Sesion.Activa, nil
}
