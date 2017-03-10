package tester

import (
	"fmt"
	"testing"
)

// We want something that have Error and Log merely
type Testing interface {
	Error(...interface{})
	Log(...interface{})
}

// ErrChecker Utils for testing
type ErrChecker struct {
	//t *testing.T
	t Testing
}

//CheckERR if err is not nil

//CheckEQ checks if two values are equals
func (ec *ErrChecker) CheckEQ(v1 interface{}, v2 interface{}) {
	if v1 != v2 {
		smsg := fmt.Sprint("\n\033[01;31mwants : ", v2, "\033[0m\n\033[01;31mgot   : ", v1, "\033[0m")
		ec.t.Error("\033[0;31mFAIL\033[0m", smsg)
		//ec.t.Error("\033[01;31m\nwants : ", v2, "\ngot   : ", v1, "\033[0m")
		return
	}
	// Only if verbose
	if testing.Verbose() {
		//ec.t.Log("\033[01;30m\nwants : ", v2, "\ngot   : ", v1, "\033[0m")
		smsg := fmt.Sprint("\n\033[01;30mwants : ", v2, "\033[0m\n\033[01;30mgot   : ", v1, "\033[0m")
		ec.t.Log("\033[00;32mPASS\033[0m", smsg)
	}
}

//MCheckEQ checks if two values are equals
func (ec *ErrChecker) MCheckEQ(msg string, v1 interface{}, v2 interface{}) {
	if v1 != v2 {
		smsg := fmt.Sprint("\n\033[01;31mwants : ", v2, "\033[0m\n\033[01;31mgot   : ", v1, "\033[0m")
		ec.t.Error("\033[0;31mFAIL\033[01;31m", msg, "\033[0m", smsg)
		//ec.t.Error("\033[01;31m\nwants : ", v2, "\ngot   : ", v1, "\033[0m")
		return
	}
	// Only if verbose
	if testing.Verbose() {
		//ec.t.Log("\033[01;30m\nwants : ", v2, "\ngot   : ", v1, "\033[0m")
		smsg := fmt.Sprint("\n\033[01;30mwants : ", v2, "\033[0m\n\033[01;30mgot   : ", v1, "\033[0m")
		ec.t.Log("\033[00;32mPASS\033[01;34m", msg, "\033[0m", smsg)
	}
}
