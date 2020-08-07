package linter

import "go/ast"

// import.go checks if there are deprecated import and sort import by group

func CheckAndFormatFile(p string) error {
	// TODO: impl, now only print path and does nothing
	log.Infof("check and format %s", p)
	return nil
}

func CheckImport(f *ast.File) error {
	return nil
}

func FormatImport(f *ast.File) error {
	return nil
}
