package database

type Repository interface {
	Get(id int64) (string, error)
	Put(id int64, value string) error
}
