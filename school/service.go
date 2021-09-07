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
	if s.ID < 0 {
		return s, errs.ErrInvalidArgument
	}
	err := svc.schoolRepo.Update(&s)
	if err != nil {
		return s, err
	}
	return s, nil
}
func (svc *schoolService) DeleteSchool(id int) (bool, error) {
	if id < 0 {
		return false, errs.ErrInvalidArgument
	}
	returnVal, err := svc.schoolRepo.Delete(id)
	return returnVal, err
}
func (svc *schoolService) GetSchools() []*School {
	schools := svc.schoolRepo.FindAll()
	return schools
}

func (svc *schoolService) GetSchool(id int) (School, error) {
	school, err := svc.schoolRepo.Find(id)
	if err != nil {
		return School{}, err
	}
	return *school, err
}
