package bg

import "fmt"

type UnexpectedStop struct {
	Service
	Err error
}

func (u *UnexpectedStop) Error() string {
	return fmt.Sprintf("service stopped unexpectedly; %s", u.Err)
}

func (u *UnexpectedStop) Unwrap() error { return u.Err }
