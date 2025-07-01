package postgres

func IsPgError(err error, target error) bool {
	return err != nil && target != nil && err.Error() == target.Error()
}