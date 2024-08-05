package version

import "github.com/tinda/cxgateway/model"

type Repository interface {
	GetServerVersion() (*model.Version, error)
	CheckVersion() error
	Upgrade() error
}
