package datastore

import (
	"bitbucket.org/enlab/peak/models"
	"bitbucket.org/enlab/peak/controllers/paginate"
)

// DataStorer interface for data layer
type DataStorer interface {
  GetSummary() (models.Summary, error)
  PutMountain(mountain *models.Mountain) (err error)
  DeleteMountain(mountainID uint64)
	GetMountains(mountainID int, page int, per_page int) (*paginate.PaginatedList, error)
	Close()
}
