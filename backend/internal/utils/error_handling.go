package utils

/*
Returns the correct huma error depending on what error is 
returned by GORM
*/
func handleGORMErrors(err error) error {
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
        	return huma.Error404NotFound("Resource not found")
    	case errors.Is(err, gorm.ErrDuplicatedKey):
        	return huma.Error409Conflict("Resource already exists", err)
   		case errors.Is(err, gorm.ErrInvalidData), errors.Is(err, gorm.ErrInvalidField), errors.Is(err, gorm.ErrInvalidValue), errors.Is(err, gorm.ErrInvalidValueOfLength):
        	return huma.Error400BadRequest("Invalid data provided", err)
    	case errors.Is(err, gorm.ErrMissingWhereClause):
        	return huma.Error400BadRequest("Missing required filter conditions")
   		case errors.Is(err, gorm.ErrPrimaryKeyRequired):
        	return huma.Error400BadRequest("Primary key is required", err)
    	case errors.Is(err, gorm.ErrInvalidTransaction):
        	return huma.Error500InternalServerError("Transaction error", err)
    	case errors.Is(err, gorm.ErrUnsupportedDriver), errors.Is(err, gorm.ErrNotImplemented):
        	return huma.Error501NotImplemented("Operation not supported", err)
		}
	}
	return nil
}

/**
Handles case when an error is returned by the DB. Ensures that nil is returned
instead of an empty entity struct
*/
func handleDBError[T any](entity *T, err error) (*T, error) {
	if err != nil {
		return nil, handleGORMErrors(err)
	}
	return entity, nil
}