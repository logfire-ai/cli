package cmdutil

type FlagError struct {
	err error
}

func FlagErrorWrap(err error) error { return &FlagError{err} }

func (fe *FlagError) Error() string {
	return fe.err.Error()
}

func (fe *FlagError) Unwrap() error {
	return fe.err
}
