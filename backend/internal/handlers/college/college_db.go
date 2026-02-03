package college

import (
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CollegeDB struct {
	db *gorm.DB
}

/*
*
Here we are using GORM to interact with the database. This is an ORM (Object Relational Mapping)
which allows us to interact with the database without having to write raw SQL queries
*/
func (c *CollegeDB) GetCollege(id uuid.UUID) (*models.College, error) {
	var college models.College
	dbResponse := c.db.Where("id = ?", id).First(&college)
	return utils.HandleDBError(&college, dbResponse.Error) // helper function that maps GORM errors to Huma errors
}

// Creates a new college in the database
func (c *CollegeDB) CreateCollege(college *models.College) (*models.College, error) {
	dbResponse := c.db.Create(college)
	return utils.HandleDBError(college, dbResponse.Error)
}

// Updates an existing college with the provided fields
func (c *CollegeDB) UpdateCollege(id uuid.UUID, update *UpdateCollegeRequest) (*models.College, error) {
	college := models.College{ID: id}

	dbResponse := c.db.
		Model(&college).
		Clauses(clause.Returning{}).
		Updates(update)

	if dbResponse.Error != nil {
		_, handleErr := utils.HandleDBError(&college, dbResponse.Error)
		return nil, handleErr
	}

	return &college, nil
}

// Performs a soft delete on a college
func (c *CollegeDB) DeleteCollege(id uuid.UUID) error {
	var college models.College

	// First check if college exists
	if err := c.db.Where("id = ?", id).First(&college).Error; err != nil {
		_, handleErr := utils.HandleDBError(&college, err)
		return handleErr
	}

	// Perform soft delete
	dbResponse := c.db.Delete(&college)
	return dbResponse.Error
}
