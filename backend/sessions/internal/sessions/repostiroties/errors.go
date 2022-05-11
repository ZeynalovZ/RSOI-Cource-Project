package repostiroties

import "errors"

var (
	NotFoundError = errors.New("no rows match specified search parameters in database")
)
