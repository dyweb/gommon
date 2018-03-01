package testutil

import (
	"testing"
	"os"
	"strings"
)

// https://github.com/dyweb/gommon/issues/57
const (
	skipMessage = "skip by gommon/testutil: %s"
	errMessage  = "gommon/testutil failed %v"
)

type Condition interface {
	Eval() (res bool, msg string, err error)
}

type con struct {
	stmt func() (res bool, msg string, err error)
}

func (c *con) Eval() (res bool, msg string, err error) {
	return c.stmt()
}

func SkipIf(t *testing.T, con Condition) {
	res, msg, err := con.Eval()
	if err != nil {
		t.Fatalf(errMessage, err)
		return
	}
	if res {
		t.Skipf(skipMessage, msg)
	} else {
		noop()
	}
}

func RunIf(t *testing.T, con Condition) {
	res, msg, err := con.Eval()
	if err != nil {
		t.Fatalf(errMessage, err)
		return
	}
	if res {
		noop()
	} else {
		t.Skipf(skipMessage, msg)
		return
	}
}

func And(l Condition, r Condition) Condition {
	return &con{
		stmt: func() (res bool, msg string, err error) {
			lr, lm, lerr := l.Eval()
			if lerr != nil {
				// TODO: wrap error instead of put it in message?
				return false, "eval And left error", lerr
			}
			if !lr {
				return false, lm, nil
			}
			rr, rm, rerr := r.Eval()
			if rerr != nil {
				return false, "eval And right error", rerr
			}
			if !rr {
				return false, rm, nil
			}
			return true, lm + " and " + rm, nil
		},
	}
}

func Or(l Condition, r Condition) Condition {
	return &con{
		stmt: func() (res bool, msg string, err error) {
			lr, lm, lerr := l.Eval()
			if lerr != nil {
				return false, "eval Or left error", lerr
			}
			if lr {
				return true, lm, nil
			}
			rr, rm, rerr := r.Eval()
			if rerr != nil {
				return false, "eval Or right error", rerr
			}
			if rr {
				return true, rm, nil
			}
			return false, lm + " or " + rm, nil
		},
	}
}

func Bool(name string, b bool) Condition {
	return &con{
		stmt: func() (res bool, msg string, err error) {
			if b {
				return true, name + " is true", nil
			} else {
				return false, name + " is false", nil
			}
		},
	}
}

func BoolP(name string, b *bool) Condition {
	return &con{
		stmt: func() (res bool, msg string, err error) {
			if *b {
				return true, name + " is true", nil
			} else {
				return false, name + " is false", nil
			}
		},
	}
}

func EnvTrue(name string) Condition {
	return &con{
		stmt: func() (res bool, msg string, err error) {
			v, ok := os.LookupEnv(name)
			if !ok {
				return false, "env " + name + " is not defined", nil
			}
			v = strings.ToLower(v)
			if v != "" && v != "0" && v != "false" && v != "f" {
				return true, "env " + name + "=true", nil
			}
			return false, "env " + name + "!=true", nil
		},
	}
}

func EnvHas(name string) Condition {
	return &con{
		stmt: func() (res bool, msg string, err error) {
			_, ok := os.LookupEnv(name)
			if ok {
				return true, "env " + name + " is defined", nil
			} else {
				return false, "env " + name + " is not defined", nil
			}
		},
	}
}

// noop does nothing
func noop() {

}
