package database

type Bucket string

const (
	Code  Bucket = "code"
	Token Bucket = "token"
)

type Repository interface {
	Get(bucket Bucket, id int64) ([]byte, error)
	Put(bucket Bucket, id int64, value []byte) error
	Delete(bucket Bucket, id int64) error
}
