package service

import (
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
