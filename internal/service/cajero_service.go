package service

import (
	"errors"

	"github.com/nicos92/dis_cajero/internal/domain/entity"
	"github.com/nicos92/dis_cajero/internal/service/ports"
)

type CajeroService struct {
	ClienteRedCentral ports.ClienteRedCentral
	ClienteSucursal   ports.ClienteSucursal
	Sesion            entity.Sesion
}

func (cajeroService *CajeroService) AuntenticarUsuario(tarjeta entity.Tarjeta, cuenta entity.Cuenta) (entity.Sesion, error) {
	validacionTarjeta, err := cajeroService.ClienteSucursal.ValidarTajeta(tarjeta)
	sesion := new(entity.Sesion)
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

func (cajeroService *CajeroService) ConsultarSaldo() (float64, entity.EstadoTransaccion, error) {
	validarSesion := cajeroService.Sesion.Activa

	if !validarSesion {
		return 0, entity.EstadoCancelado, errors.New("La cuenta sesion no esta activa")
	}

	saldo, err := cajeroService.ClienteRedCentral.ConsultarSaldo(cajeroService.Sesion.CuentaCliente)

	if err != nil {
		return 0, entity.EstadoCancelado, err
	}

	return saldo, entity.EstadoConfirmado, nil
}

func (cajeroService *CajeroService) ExtraerEfectivo(monto float64) (entity.EstadoTransaccion, int, error) {

	validarSesion := cajeroService.Sesion.Activa

	if !validarSesion {
		return entity.EstadoCancelado, 0, errors.New("La cuenta sesion no esta activa")
	}

	if cajeroService.Sesion.CuentaCliente.Saldo < monto {
		return entity.EstadoCancelado, 0, errors.New("Saldo insuficiente en la cuenta.")
	}
	estado, id, err := cajeroService.ClienteRedCentral.DebitarExtaccion(cajeroService.Sesion.CuentaCliente, monto)

	if err != nil {
		return estado, id, err
	}

	cajeroService.Sesion.CuentaCliente.Saldo -= monto

	return estado, id, nil
}

func (cajeroService *CajeroService) RealizarTransferencia(cuentaDestino entity.Cuenta, monto float64) (entity.EstadoTransaccion, int, error) {
	validarSesion := cajeroService.Sesion.Activa

	if !validarSesion {
		return entity.EstadoCancelado, 0, errors.New("La cuenta sesion no esta activa")
	}

	if cajeroService.Sesion.CuentaCliente.Saldo < monto {
		return entity.EstadoCancelado, 0, errors.New("Saldo insuficiente en la cuenta.")
	}

	estado, id, err := cajeroService.ClienteRedCentral.TransferirACuenta(cajeroService.Sesion.CuentaCliente, cuentaDestino, monto)

	if err != nil {
		return estado, id, err
	}

	cajeroService.Sesion.CuentaCliente.Saldo -= monto

	return estado, id, nil
}
