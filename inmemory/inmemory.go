package inmemory

import (
	"sync"

	schools "schoolmgt/school"

	errs "schoolmgt/errors"
)

// inMemUserRepository is an implementation of a school repository for storage in local memory
type inMemUserRepository struct {
	mtx     *sync.RWMutex
	schools map[int]*schools.School
}

// NewInMemUserRepository returns a new school repository for storage in local memory
func NewInMemUserRepository() schools.Repository {
	return &inMemUserRepository{
		mtx:     new(sync.RWMutex),
		schools: make(map[int]*schools.School),
	}
}

// Store inserts a school into the local school map
func (ir *inMemUserRepository) Store(school *schools.School) error {
	ir.mtx.Lock()
	school.ID = len(ir.schools) + 1
	ir.schools[school.ID] = school
	ir.mtx.Unlock()
	return nil
}

// Updates a school into the local school map
func (ir *inMemUserRepository) Update(school *schools.School) error {
	ir.mtx.Lock()
	var err error
	_, ok := ir.schools[school.ID]
	if ok {
		ir.schools[school.ID] = school
		err = nil
	} else {
		err = errs.ErrSchoolNotFound
	}
	ir.mtx.Unlock()
	return err
}

// Find retrieves a single school from the repository
func (ir *inMemUserRepository) Find(id int) (*schools.School, error) {
	ir.mtx.RLock()
	s, ok := ir.schools[id]
	ir.mtx.RUnlock()
	if !ok {
		return nil, errs.ErrSchoolNotFound
	}
	return s, nil
}

// FindAll retrieves all schools from memory
func (ir *inMemUserRepository) FindAll() []*schools.School {
	ir.mtx.RLock()
	allUsers := []*schools.School{}
	for _, v := range ir.schools {
		allUsers = append(allUsers, v)
	}
	ir.mtx.RUnlock()
	return allUsers
}

func (ir *inMemUserRepository) Delete(id int) (bool, error) {
	ir.mtx.Lock()
	_, ok := ir.schools[id]
	if ok {
		delete(ir.schools, id)
		ir.mtx.Unlock()
	} else {
		ir.mtx.Unlock()
		return false, errs.ErrSchoolNotFound
	}
	return true, nil
}
