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

type EstadoTransaccion string

const (
	EstadoProceso    EstadoTransaccion = "Procesando"
	EstadoCancelado  EstadoTransaccion = "Cancelada"
	EstadoConfirmado EstadoTransaccion = "Confirmada"
)

type TipoTransaccion string

const (
	Consulta      TipoTransaccion = "Consulta"
	Extraccion    TipoTransaccion = "Extraccion"
	Transferencia TipoTransaccion = "Transferencia"
)
