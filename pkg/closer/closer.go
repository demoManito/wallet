package closer

import "io"

type CloseFunc func() error

func NewCloserDelegate(f func() error) io.Closer {
	if f == nil {
		return nil
	}
	return &CloserDelegate{
		Func: f,
	}
}

type CloserDelegate struct {
	Func func() error
}

func (c *CloserDelegate) Close() error {
	return c.Func()
}

func NopClose() error { return nil }

func Pocket(interruptable bool, closeFunc ...CloseFunc) CloseFunc {
	return func() error {
		var err error
		for _, f := range closeFunc {
			err = f()
			if interruptable && err != nil {
				return err
			}
		}
		return nil
	}
}
