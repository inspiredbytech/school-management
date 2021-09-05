package stringsvc

type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
}
