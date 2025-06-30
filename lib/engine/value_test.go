package engine

type Test struct {
	value string
}

func (test Test) Headers() []string {
	return nil
}

func (test Test) Values() []string {
	return nil
}

func (test Test) Key() string   { return "" }
func (test Test) Scope() string { return "" }
func (test Test) Type() string  { return "tuType" }
