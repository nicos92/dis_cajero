package entidad

type Sesion struct {
	Tarjeta       Tarjeta
	CuentaCliente Cuenta
	Activa        bool
}
