package ports

import (
	"github.com/nicos92/dis_cajero/internal/domain/entity"
)

type ClienteSucursal interface {
	ValidarTajeta(tarjeta entity.Tarjeta) (bool, error)
}
