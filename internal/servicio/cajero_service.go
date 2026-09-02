package servicio

import (
	"errors"

	entidad "github.com/nicos92/dis_cajero/internal/dominio/entidad"
	"github.com/nicos92/dis_cajero/internal/servicio/ports"
)

type CajeroService struct {
	ClienteRedCentral ports.ClienteRedCentral
	ClienteSucursal   ports.ClienteSucursal
	Sesion            entidad.Sesion
}

func (cajeroService *CajeroService) AuntenticarUsuario(tarjeta entidad.Tarjeta, cuenta entidad.Cuenta) (entidad.Sesion, error) {
	validacionTarjeta, err := cajeroService.ClienteSucursal.ValidarTarjeta(tarjeta)
	sesion := new(entidad.Sesion)
	sesion.Activa = false
	if err != nil {
		return *sesion, err
	}

	validacionPin, err := cajeroService.ClienteRedCentral.ValidarPin(cuenta)

	if err != nil {
		return *sesion, err
	}

	if !validacionTarjeta || !validacionPin {
		return *sesion, err
	}

	sesion.Activa = true
	sesion.CuentaCliente = cuenta
	cajeroService.Sesion = *sesion
	return *sesion, nil

}

func (cajeroService *CajeroService) ConsultarSaldo() (float64, entidad.EstadoTransaccion, error) {
	validarSesion := cajeroService.Sesion.Activa

	if !validarSesion {
		return 0, entidad.EstadoCancelado, errors.New("La cuenta sesion no esta activa")
	}

	saldo, err := cajeroService.ClienteRedCentral.ConsultarSaldo(cajeroService.Sesion.CuentaCliente)

	if err != nil {
		return 0, entidad.EstadoCancelado, err
	}

	return saldo, entidad.EstadoConfirmado, nil
}

func (cajeroService *CajeroService) ExtraerEfectivo(monto float64) (entidad.EstadoTransaccion, int, error) {

	validarSesion := cajeroService.Sesion.Activa

	if !validarSesion {
		return entidad.EstadoCancelado, 0, errors.New("La cuenta sesion no esta activa")
	}

	if cajeroService.Sesion.CuentaCliente.Saldo < monto {
		return entidad.EstadoCancelado, 0, errors.New("Saldo insuficiente en la cuenta.")
	}
	estado, id, err := cajeroService.ClienteRedCentral.DebitarExtraccion(cajeroService.Sesion.CuentaCliente, monto)

	if err != nil {
		return estado, id, err
	}

	cajeroService.Sesion.CuentaCliente.Saldo -= monto

	return estado, id, nil
}

func (cajeroService *CajeroService) RealizarTransferencia(cuentaDestino entidad.Cuenta, monto float64) (entidad.EstadoTransaccion, int, error) {
	validarSesion := cajeroService.Sesion.Activa

	if !validarSesion {
		return entidad.EstadoCancelado, 0, errors.New("La cuenta sesion no esta activa")
	}

	if cajeroService.Sesion.CuentaCliente.Saldo < monto {
		return entidad.EstadoCancelado, 0, errors.New("Saldo insuficiente en la cuenta.")
	}

	estado, id, err := cajeroService.ClienteRedCentral.TransferirACuenta(cajeroService.Sesion.CuentaCliente, cuentaDestino, monto)

	if err != nil {
		return estado, id, err
	}

	cajeroService.Sesion.CuentaCliente.Saldo -= monto

	return estado, id, nil
}

func (cajeroService *CajeroService) CerrarSesion() (bool, error) {
	validarSesion := cajeroService.Sesion.Activa

	if !validarSesion {
		return false, errors.New("La cuenta sesion no esta activa")
	}

	cajeroService.Sesion.Activa = false
	return cajeroService.Sesion.Activa, nil
}
