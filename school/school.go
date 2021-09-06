package school

type School struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Country  string   `json:"country"`
	City     string   `json:"city"`
	Address  string   `json:"Address"`
	Contacts []string `json:"Contacts"`
}

// New returns a reference to a new school instance
func New(id int, name string, country string, city string, address string, contacts []string) *School {
	return &School{
		ID:       id,
		Name:     name,
		Country:  country,
		City:     city,
		Address:  address,
		Contacts: contacts,
	}
}

// (School) Repository is the set of behavior a repository, or "store", of schools must conform to.
type Repository interface {
	// Store a new user in the repository
	Store(user *School) error

	// Find a user in the repository by ID
	Find(id int) (*School, error)

	// FindAll users in the repository
	FindAll() []*School
}
