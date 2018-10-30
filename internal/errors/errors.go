package errors

func New(code, msg string) error {
	return Error{Code: code, Msg: msg}
}

type Error struct {
	Code string
	Msg  string
}

func (e Error) Error() string {
	return e.Msg
}
