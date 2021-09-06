package school

import errs "schoolmgt/errors"

//Service describes the behavior of the school service, CRUD Activities
type Service interface {
	CreateSchool(s School) (int, error)
	UpdateSchool(s School) (School, error)
	DeleteSchool(id int) (bool, error)
	GetSchools() []*School
	GetSchool(id int) (School, error)
}

//implementation of the school service
type schoolService struct {
	schoolRepo Repository
}

func NewService(repo Repository) Service {
	return &schoolService{
		schoolRepo: repo,
	}
}

func (svc *schoolService) CreateSchool(s School) (int, error) {
	if s.ID < 0 {
		return s.ID, errs.ErrInvalidArgument
	}
	err := svc.schoolRepo.Store(&s)
	if err != nil {
		return s.ID, err
	}
	return s.ID, nil
}

func (svc *schoolService) UpdateSchool(s School) (School, error) {
	return School{}, nil
}
func (svc *schoolService) DeleteSchool(id int) (bool, error) {
	return true, nil
}
func (svc *schoolService) GetSchools() []*School {
	schools := svc.schoolRepo.FindAll()

	return schools
}

func (svc *schoolService) GetSchool(id int) (School, error) {
	return School{}, nil
}
