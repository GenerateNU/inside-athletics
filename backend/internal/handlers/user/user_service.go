package user

import (
	"context"
	"inside-athletics/internal/handlers/role"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/s3"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type UserService struct {
	userDB *UserDB
	roleDB *role.RoleDB
	s3     *s3.Service
}

/*
*
Huma automatically validates the params based on which type you have passed it based on the struct tags (so nice!!)
IF there is an error it will automatically send the correct response back with the right status and message about the validation errors

Here we are mapping to a GetUserResponse so that we can control what the return type is. It's important to make seperate
return types so that we can control what information we are sending back instead of just the entire model
*/
func (u *UserService) GetUser(ctx context.Context, input *GetUserParams) (*utils.ResponseBody[GetUserResponse], error) {
	id := input.ID
	user, err := u.userDB.GetUser(id)
	respBody := &utils.ResponseBody[GetUserResponse]{}

	if err != nil {
		return respBody, err
	}

	roleResponses, err := u.userDB.GetRolesWithPermissionsForUser(id)
	if err != nil {
		return nil, err
	}
	response := u.buildUserResponse(ctx, user, roleResponses)

	return &utils.ResponseBody[GetUserResponse]{
		Body: response,
	}, err
}

func (u *UserService) GetCurrentUser(ctx context.Context, input *utils.EmptyInput) (*utils.ResponseBody[GetUserResponse], error) {
	respBody := &utils.ResponseBody[GetUserResponse]{}

	userID, err := u.getCurrentUserID(ctx)
	if err != nil {
		return respBody, err
	}

	user, err := u.userDB.GetUser(userID)
	if err != nil {
		return respBody, err
	}

	roleResponses, err := u.userDB.GetRolesWithPermissionsForUser(userID)
	if err != nil {
		return nil, err
	}

	respBody.Body = u.buildUserResponse(ctx, user, roleResponses)

	return respBody, nil
}

func (u *UserService) CreateUser(ctx context.Context, input *CreateUserInput) (*utils.ResponseBody[CreateUserResponse], error) {
	respBody := &utils.ResponseBody[CreateUserResponse]{}

	currentUserID, err := u.getCurrentUserID(ctx)
	if err != nil {
		return respBody, err
	}

	roleID, err := u.roleDB.GetRoleIDByName(models.RoleUser)
	if err != nil {
		return respBody, err
	}

	user := &models.User{
		ID:                      currentUserID,
		FirstName:               input.Body.FirstName,
		LastName:                input.Body.LastName,
		Email:                   input.Body.Email,
		Username:                input.Body.Username,
		Bio:                     input.Body.Bio,
		ProfileImageKey:         input.Body.ProfileImageKey,
		Account_Type:            input.Body.AccountType,
		SportID:                 input.Body.SportID,
		Expected_Grad_Year:      input.Body.ExpectedGradYear,
		Verified_Athlete_Status: input.Body.VerifiedAthleteStatus,
		CollegeID:               input.Body.CollegeID,
		Division:                input.Body.Division,
	}

	createdUser, err := u.userDB.CreateUser(user)
	if err != nil {
		return respBody, err
	}

	if err := u.userDB.AddUserRole(createdUser.ID, roleID); err != nil {
		return respBody, err
	}

	respBody.Body = &CreateUserResponse{
		ID:   createdUser.ID,
		Name: createdUser.FirstName,
		Role: &UserRoleResponse{
			ID:   roleID,
			Name: models.RoleUser,
		},
	}

	return respBody, nil
}

func (u *UserService) UpdateUser(ctx context.Context, input *UpdateUserInput) (*utils.ResponseBody[UpdateUserResponse], error) {
	respBody := &utils.ResponseBody[UpdateUserResponse]{}
	currentUserID, err := u.getCurrentUserID(ctx)
	if err != nil {
		return respBody, err
	}

	updatedUser, err := u.userDB.UpdateUser(currentUserID, input.Body)
	if err != nil {
		return respBody, err
	}

	roleResponses, err := u.userDB.GetRolesWithPermissionsForUser(currentUserID)
	if err != nil {
		return nil, err
	}

	respBody.Body = u.buildUserResponse(ctx, updatedUser, roleResponses)

	return respBody, nil
}

func (u *UserService) DeleteUser(ctx context.Context, input *GetUserParams) (*utils.ResponseBody[DeleteUserResponse], error) {
	respBody := &utils.ResponseBody[DeleteUserResponse]{}

	err := u.userDB.DeleteUser(input.ID)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &DeleteUserResponse{
		ID: input.ID,
	}

	return respBody, nil
}

func (u *UserService) AssignRole(ctx context.Context, input *AssignRoleInput) (*utils.ResponseBody[AssignRoleResponse], error) {
	if input.Body.RoleID == uuid.Nil {
		return nil, huma.Error422UnprocessableEntity("role_id cannot be empty")
	}

	_, err := u.userDB.GetUser(input.ID)
	if err != nil {
		return nil, err
	}

	role, err := u.roleDB.GetRoleByID(input.Body.RoleID)
	if err != nil {
		return nil, err
	}

	if err := u.userDB.AddUserRole(input.ID, input.Body.RoleID); err != nil {
		return nil, err
	}

	return &utils.ResponseBody[AssignRoleResponse]{
		Body: &AssignRoleResponse{
			UserID: input.ID,
			Role: UserRoleResponse{
				ID:   role.ID,
				Name: role.Name,
			},
		},
	}, nil
}

func (u *UserService) getCurrentUserID(ctx context.Context) (uuid.UUID, error) {
	rawID := ctx.Value("user_id")
	if rawID == nil {
		return uuid.Nil, huma.Error401Unauthorized("User not authenticated")
	}

	userID, ok := rawID.(string)
	if !ok {
		return uuid.Nil, huma.Error500InternalServerError("Invalid user ID in context")
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, huma.Error400BadRequest("Invalid user ID", err)
	}

	return parsedID, nil
}

func (u *UserService) buildUserResponse(ctx context.Context, user *models.User, roleResponses *[]role.RoleResponse) *GetUserResponse {
	response := &GetUserResponse{
		ID:                    user.ID,
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		Email:                 user.Email,
		Username:              user.Username,
		Bio:                   user.Bio,
		ProfileImageKey:       user.ProfileImageKey,
		AccountType:           user.Account_Type,
		Sport:                 user.Sport,
		ExpectedGradYear:      user.Expected_Grad_Year,
		VerifiedAthleteStatus: user.Verified_Athlete_Status,
		College:               user.College,
		Division:              user.Division,
		Roles:                 roleResponses,
	}

	if u.s3 != nil && user.ProfileImageKey != nil && *user.ProfileImageKey != "" {
		if download, err := u.s3.GetDownloadURL(ctx, *user.ProfileImageKey); err == nil {
			response.ProfileImageURL = &download.DownloadURL
		}
	}

	return response
}
