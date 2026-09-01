package almacenamiento

import entidad "github.com/nicos92/dis_cajero/internal/dominio/entidad"

type RedCentralRepositorio struct {
	Cuentas map[string]entidad.Cuenta
}
