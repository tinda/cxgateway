package version

import "github.com/tinda/cxgateway/model"

type Service interface {
	GetServerVersion() (*model.Version, error)
	CheckVersion() error
	Upgrade() error
}
