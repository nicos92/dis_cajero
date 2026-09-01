package entidad

import "time"

type Transaccion struct {
	Id            int
	Tipo          TipoTransaccion
	CuentaOrigen  Cuenta
	CuentaDestino Cuenta
	Monto         float64
	Creado        time.Time
	Estado        EstadoTransaccion
}

type EstadoTransaccion int

const (
	EstadoProceso EstadoTransaccion = iota
	EstadoCancelado
	EstadoConfirmado
)

type TipoTransaccion string

const (
	Consulta      = "Consulta"
	Extraccion    = "Extraccion"
	Transferencia = "Transferencia"
)
