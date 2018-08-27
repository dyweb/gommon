// Package genutil contains helper when generating files, it is used break dependency cycle between generator package
// and packages that contain generator logic like log, noodle
package genutil

// DefaultHeader returns the standard go header for generated files with two trailing \n,
// the second \n is to avoid this header becomes documentation of the package
func DefaultHeader(templateSrc string) string {
	return "// Code generated by gommon from " + templateSrc + " DO NOT EDIT.\n\n"
}

// Header returns the standard go header for generated files with two trailing \n,
// the second \n is to avoid this header becomes documentation of the package
func Header(generator string, templateSrc string) string {
	return "// Code generated by " + generator + " from " + templateSrc + " DO NOT EDIT.\n\n"
}
