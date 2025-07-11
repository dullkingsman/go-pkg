package prizzle

import "github.com/lib/pq"

type DbError pq.Error

type DatabaseError[T error] struct {
	DbError     pq.Error
	driver      string
	NativeError T
}

func (e *DatabaseError[T]) Error() string {
	return e.DbError.Message
}
