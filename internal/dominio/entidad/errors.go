package entidad

import "errors"

var (
	ErrSesionInactiva        = errors.New("sesión no activa")
	ErrSaldoInsuficiente     = errors.New("saldo insuficiente")
	ErrCredencialesInvalidas = errors.New("credenciales inválidas")
	ErrPinIncorrecto         = errors.New("pin incorrecto")
	ErrCuentaNoExiste        = errors.New("cuenta no existe")
	ErrCuentaInactiva        = errors.New("cuenta inactiva")
	ErrCuentaDestinoInvalida = errors.New("cuenta de destino inválida")
	ErrCuentasNoAsociadas    = errors.New("la tarjeta no tiene cuentas asociadas")
	ErrTarjetaInvalida       = errors.New("tarjeta inválida")
	ErrTarjetaVencida        = errors.New("tarjeta vencida")
)
