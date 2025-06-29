package postgres

import "errors"

var (
	ErrUsersEmailDuplicate = errors.New("pq: duplicate key value violates unique constraint \"users_email_key\"")
	ErrUsersNumberDuplicate = errors.New("pq: duplicate key value violates unique constraint \"users_number_key\"")
)