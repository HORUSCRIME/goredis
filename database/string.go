package database

type String struct {
	Val string
}

func NewString(val string) *String {
	return &String{Val: val}
}

func (s *String) Type() string {
	return "string"
}