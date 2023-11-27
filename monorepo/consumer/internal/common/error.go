package common

type WrapError struct {
	Code int
	Err  error
	Msg  string
	Body any
}

func (w WrapError) Error() string {
	return w.Err.Error()
}

func (w WrapError) Unwrap() error {
	return w.Err
}

func (w WrapError) Message() string {
	return w.Msg
}
