package testutil

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// https://github.com/dyweb/gommon/issues/57
const (
	skipMessage = "skip by gommon/testutil: %s"
	errMessage  = "gommon/testutil failed %v"
)

type Condition interface {
	Eval() (res bool, msg string, err error)
	// B calls Eval and returns false if there is an error when eval.
	B() bool
}

type con struct {
	stmt func() (res bool, msg string, err error)
}

func (c *con) Eval() (res bool, msg string, err error) {
	return c.stmt()
}

func (c *con) B() bool {
	b, _, err := c.stmt()
	if err != nil || b == false {
		return false
	}
	return true
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

// Start of boolean expressions, And, Or, Not

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

func Not(c Condition) Condition {
	return &con{
		stmt: func() (res bool, msg string, err error) {
			r, m, e := c.Eval()
			if e != nil {
				return false, "eval Not error", e
			}
			return !r, "not " + m, e
		},
	}
}

// end of boolean expressions.

// Bool checks if a bool is true
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

// BoolP checks if value of a pointer to bool is true
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

// BinaryExist returns a test condition that looks up binary from PATH using exec.LookPath.
func BinaryExist(name string) Condition {
	return &con{stmt: func() (res bool, msg string, err error) {
		p, err := exec.LookPath(name)
		if err == nil {
			return true, fmt.Sprintf("binary %s is in %s", name, p), nil
		}
		return false, fmt.Sprintf("binary %s not found: %s", name, err), err
	}}
}

// CommandSuccess returns a test condition that executes a command and wait its completion.
// It evaluates to true if the process exist with 0.
// TODO: maybe have a default timeout?
func CommandSuccess(cmd string, args ...string) Condition {
	return &con{stmt: func() (res bool, msg string, err error) {
		cmd := exec.Command(cmd, args...)
		full := fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
		b, err := cmd.CombinedOutput()
		if err == nil {
			return true, "command success: " + full, nil
		}
		// FIXME: we return nil error to avoid failing the test in RunIf, might need to change the design
		return false, "command failed: " + full + string(b), nil
	}}
}

// wrapper for common conditions, NOTE: some are defined in other files like golden.go

// IsTravis check if env TRAVIS is true
func IsTravis() Condition {
	return EnvTrue("TRAVIS")
}

// HasDocker returns a test condition that checks if docker client exists and the daemon is up and running.
// TODO: cache the result of docker version?
func HasDocker() Condition {
	return And(BinaryExist("docker"), CommandSuccess("docker", "version"))
}

// Dump check if env DUMP or GOMMON_DUMP is set, so print detail or use go-spew to dump structs etc.
func Dump() Condition {
	return Or(EnvHas("DUMP"), EnvHas("GOMMON_DUMP"))
}

// noop does nothing
func noop() {

}
