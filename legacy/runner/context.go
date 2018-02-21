package runner

// FIXME: this is never used, and it may not be a good idea to use context, since there is context package
// Context describes how the command should be run
// i.e. use shell instead of using os/exec,
// use self defined logic instead of lookup an executable
// run in background
type Context struct {
	AutoShell  bool
	Foreground bool
	Block      bool
	Dry        bool
}

// NewContext returns a context with convention
func NewContext() *Context {
	return &Context{
		AutoShell:  true,
		Foreground: true,
		Block:      true,
		Dry:        false,
	}
}
