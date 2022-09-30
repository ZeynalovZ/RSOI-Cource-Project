package repositories

import "errors"

var (
	NotFoundError = errors.New("no rows match specified search parameters in database")
	PaymentAlreadyExists = errors.New("payments count in database not equals zero")
)
