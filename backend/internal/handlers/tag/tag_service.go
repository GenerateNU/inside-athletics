package tag

import (
	"context"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"reflect"
	"strings"
)

type TagService struct {
	tagDB *TagDB
}

func (u *TagService) GetTagByName(ctx context.Context, input *GetTagByNameParams) (*utils.ResponseBody[GetTagResponse], error) {
	name := input.Name
	tag, err := u.tagDB.GetTagByName(name)
	respBody := &utils.ResponseBody[GetTagResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetTagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}

	return &utils.ResponseBody[GetTagResponse]{
		Body: response,
	}, err
}

func (u *TagService) GetTagById(ctx context.Context, input *GetTagByIDParams) (*utils.ResponseBody[GetTagResponse], error) {
	id := input.ID
	tag, err := u.tagDB.GetTagByID(id)
	respBody := &utils.ResponseBody[GetTagResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetTagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}

	return &utils.ResponseBody[GetTagResponse]{
		Body: response,
	}, err
}

func (u *TagService) CreateTag(ctx context.Context, input *CreateTagInput) (*utils.ResponseBody[CreateTagResponse], error) {
	respBody := &utils.ResponseBody[CreateTagResponse]{}

	tag := &models.Tag{
		Name: input.Body.Name,
	}

	createdTag, err := u.tagDB.CreateTag(tag)

	if err != nil {
		return respBody, err
	}

	response := &CreateTagResponse{
		ID:   createdTag.ID,
		Name: createdTag.Name,
	}

	return &utils.ResponseBody[CreateTagResponse]{
		Body: response,
	}, err
}

func (u *TagService) UpdateTag(cts context.Context, input *UpdateTagInput) (*utils.ResponseBody[UpdateTagResponse], error) {
	respBody := &utils.ResponseBody[UpdateTagResponse]{}

	updates, err := buildTagUpdates(input.Body)
	if err != nil {
		return respBody, err
	}

	updatedTag, err := u.tagDB.UpdateTag(input.ID, updates)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &UpdateTagResponse{
		ID:   updatedTag.ID,
		Name: updatedTag.Name,
	}

	return respBody, nil
}

func buildTagUpdates(body UpdateTagBody) (map[string]interface{}, error) {
	updates := make(map[string]interface{})
	val := reflect.ValueOf(body)
	typ := reflect.TypeOf(body)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		name := strings.Split(tag, ",")[0]
		if name == "" {
			continue
		}

		fieldVal := val.Field(i)
		if fieldVal.IsNil() {
			continue
		}

		updates[name] = fieldVal.Elem().Interface()
	}

	return updates, nil
}

func (u *TagService) DeleteTag(ctx context.Context, input *GetTagByIDParams) (*utils.ResponseBody[DeleteTagResponse], error) {
	respBody := &utils.ResponseBody[DeleteTagResponse]{}

	err := u.tagDB.DeleteTag(input.ID)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &DeleteTagResponse{
		ID: input.ID,
	}

	return respBody, nil
}
