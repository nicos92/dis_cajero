package entidad

import "errors"

var (
	ErrSesionInactiva        = errors.New("sesión no activa")
	ErrSaldoInsuficiente     = errors.New("saldo insuficiente")
	ErrCredencialesInvalidas = errors.New("credenciales inválidas")
	ErrCuentaNoExiste        = errors.New("cuenta no existe")
	ErrTarjetaInvalida       = errors.New("tarjeta inválida")
	ErrTarjetaVencida        = errors.New("tarjeta vencida")
)
