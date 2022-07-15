package tests

type Test struct {
	Name     string
	Function func() (bool, error)
}
