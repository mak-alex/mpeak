package datastore

import (
	_ "os"

	"bitbucket.org/enlab/peak/models"
	"bitbucket.org/enlab/peak/controllers/paginate"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	DEFAULT_PAGE_SIZE int = 20
	MAX_PAGE_SIZE     int = 1000
)

type dbStore struct {
	db *gorm.DB
}

func (store *dbStore) GetMountains(mountainID int, page int, per_page int) (*paginate.PaginatedList, error) {
	count := 0
	result := []models.Mountain{}

	m := store.db.Select("mountains*").Table("mountains")
  if mountainID != 0 {
    //for _, term := range utils.SplitBySeparators(strings.ToLower(mountain)) {
    m = m.Where("id = ?", mountainID)
  }
	if per_page <= 0 {
		per_page = DEFAULT_PAGE_SIZE
	}
	if per_page > MAX_PAGE_SIZE {
		per_page = MAX_PAGE_SIZE
	}
	if page == 0 {
		page = 1
	}
	m.Count(&count)
	p := paginate.NewPaginatedList(page, per_page, count)
	m = m.Limit(p.Limit())
	m = m.Offset(p.Offset())
	m.Select("mountains.*").Table("mountains").Find(&result)

	p.Items = result
	return p, nil
}

func (store *dbStore) PutMountain(mountain *models.Mountain) (err error) {
	tx := store.db.Begin()
	store.db.Create(&mountain)
	tx.Commit()

	return err
}

func (store *dbStore) DeleteMountain(mountainID uint64) {
  store.db.Where("id = ?", mountainID).Delete(&models.Mountain{})
}

func (store *dbStore) GetSummary() (models.Summary, error) {
	summary := models.Summary{}
	mountains_count := 0
	store.db.Table("mountains").Count(&mountains_count)
	summary.Mountains = mountains_count

	return summary, nil
}

func (store *dbStore) Close() {
	store.db.Close()
}

// NewDBStore creates new instance of datastorer
func NewDBStore(config *models.DBConfig) (DataStorer, error) {
	db, err := gorm.Open(config.DBType, config.DBParams)
	if err == nil {
		db.DB()
		db.AutoMigrate(&models.Mountain{})
		db.LogMode(config.DBLog)
	}
	result := new(dbStore)
	result.db = db

	return result, err
}
