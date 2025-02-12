package repositories

import "errors"

var (
	ErrRecordNotFound        = errors.New("record was not found")
	ErrDatabaseReadingError  = errors.New("error while reading from DB")
	ErrRecordAlreadyExists   = errors.New("record already exists")
	ErrDatabaseUpdatingError = errors.New("record was not updated")
	ErrStartTxError          = errors.New("failed to start transaction")
	ErrRollbackTxError       = errors.New("failed to rollback transaction")
	ErrCommitError           = errors.New("failed to commit transaction")
)
