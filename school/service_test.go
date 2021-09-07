package school

import (
	"testing"

	"errors"

	errs "schoolmgt/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSchoolService_CreateSchoolInvalidArgs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	us := NewService(mockRepo)
	mockSchool := NewMockSchool()
	id, err := us.CreateSchool(mockSchool)
	assert.Equal(t, mockSchool.ID, id)
	assert.EqualError(t, err, errs.ErrInvalidArgument.Error())
}

func TestSchoolService_CreateSchoolFailureOnStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	us := NewService(mockRepo)
	mockSchool := NewMockSchool()
	mockRepo.EXPECT().Store(&mockSchool).Return(errors.New("I'm a repository error!"))
	id, err := us.CreateSchool(mockSchool)
	assert.Error(t, err)
	assert.Equal(t, mockSchool.ID, id)
}

func TestSchoolService_CreateSchoolSuccessfulStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockSchool := NewMockSchool()
	us := NewService(mockRepo)
	mockRepo.EXPECT().Store(&mockSchool).Return(nil)
	id, err := us.CreateSchool(mockSchool)
	assert.NoError(t, err)
	assert.Equal(t, mockSchool.ID, id)
}

func NewMockSchool() School {
	return School{
		ID:       1,
		Name:     "St John's",
		Country:  "UK",
		City:     "London",
		Address:  "",
		Contacts: []string{"Mahesh"},
	}
}
