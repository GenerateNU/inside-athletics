package utils


/*
Returns the correct huma error depending on what error is 
returned by GORM
*/
func handleDbErrors(dbResponse *gorm.DB) error {
	if dbResponse.error != nil {
		switch {
		case errors.Is(dbResponse.error, gorm.ErrRecordNotFound):
        	return huma.Error404NotFound("Resource not found")
    	case errors.Is(dbResponse.error, gorm.ErrDuplicatedKey):
        	return huma.Error409Conflict("Resource already exists", dbResponse.error)
   		case errors.Is(dbResponse.error, gorm.ErrInvalidData), errors.Is(err, gorm.ErrInvalidField), errors.Is(err, gorm.ErrInvalidValue), errors.Is(err, gorm.ErrInvalidValueOfLength):
        	return huma.Error400BadRequest("Invalid data provided", dbResponse.error)
    	case errors.Is(dbResponse.error, gorm.ErrMissingWhereClause):
        	return huma.Error400BadRequest("Missing required filter conditions")
   		case errors.Is(dbResponse.error, gorm.ErrPrimaryKeyRequired):
        	return huma.Error400BadRequest("Primary key is required", dbResponse.error)
    	case errors.Is(dbResponse.error, gorm.ErrInvalidTransaction):
        	return huma.Error500InternalServerError("Transaction error", dbResponse.error)
    	case errors.Is(dbResponse.error, gorm.ErrUnsupportedDriver), errors.Is(err, gorm.ErrNotImplemented):
        	return huma.Error501NotImplemented("Operation not supported", dbResponse.error)
		}
	}
	return nil
}