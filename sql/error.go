package sql

import "fmt"

type multiErr struct {
	Errors []error
}

func (me multiErr) Error() string {
	return fmt.Sprint(me.Errors)
}

func (me multiErr) orNil() error {
	if len(me.Errors) == 0 {
		return nil
	}
	return me
}

// append returns err, not me.
func (me *multiErr) append(err error) error {
	if err != nil {
		me.Errors = append(me.Errors, err)
	}
	return err
}
