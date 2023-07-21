package database

type Bucket string

const (
	Code  Bucket = "code"
	Token Bucket = "token"
)

type Repository interface {
	Get(bucket Bucket, id int64) (string, error)
	Put(bucket Bucket, id int64, value string) error
	Delete(bucket Bucket, id int64) error
}
