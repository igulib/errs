package errs

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type ErrorDetails struct {
	Show         ElementsToShow
	PkgFullName  string
	PkgShortName string
}

type ElementsToShow struct {
	// Show package name
	Pkg bool
	// Show function name
	Func bool
	// Show file and line
	File bool
	// Has effect only when Pkg is true
	FullPackageName bool
	// Has effect only when File is true
	FilePath bool
}

// Sep is a separator string to separate wrapped errors when creating
// a string representation of the error.
var Sep string = ": ->\n    <- "

func NewErrorDetails(pkgFullName string) *ErrorDetails {
	pkgShortName := filepath.Base(pkgFullName)
	return &ErrorDetails{
		Show: ElementsToShow{
			Pkg:  true,
			Func: true,
			File: true,
		},
		PkgFullName:  pkgFullName,
		PkgShortName: pkgShortName,
	}
}

func Format(v *ErrorDetails, indirectionSteps int, msgFmt string, args ...interface{}) string {
	var (
		pkg, fn, fileAndLine string
	)

	// Handle case when ErrorDetails is nil
	if v == nil {
		pc, file, lineNum, ok := runtime.Caller(indirectionSteps)
		if ok {
			// Get caller function name
			fn = runtime.FuncForPC(pc).Name()
			// Discard package name
			dotPos := strings.LastIndex(fn, ".")
			if dotPos != -1 {
				fn = fn[dotPos+1:]
			}
			file = filepath.Base(file)
			fileAndLine = fmt.Sprintf(", %s:%d", file, lineNum)
		}
		if msgFmt != "" {
			if len(args) > 0 {
				msgFmt = fmt.Sprintf(msgFmt, args...)
			}
			return fmt.Sprintf("[!ErrorDetails==nil] in %s%s (%s)",
				fn, fileAndLine, msgFmt)
		}

		return fmt.Sprintf("[!ErrorDetails==nil] in %s%s", fn, fileAndLine)
	}

	var (
		showPkg  = v.Show.Pkg
		showFunc = v.Show.Func
		showFile = v.Show.File
	)

	if showPkg {
		if !v.Show.FullPackageName {
			pkg = "pkg " + v.PkgShortName
		} else {
			pkg = "pkg " + v.PkgFullName
		}
	}

	if showFunc || showFile {
		pc, file, lineNum, ok := runtime.Caller(indirectionSteps)
		if ok {
			if showFunc {
				fn = runtime.FuncForPC(pc).Name()
				// Get only function name
				dotPos := strings.LastIndex(fn, ".")
				if dotPos != -1 && dotPos < len(fn) {
					if showPkg {
						// Include dot only if pkg name is shown
						fn = fn[dotPos:]
					} else {
						fn = fn[dotPos+1:]
					}
				}
			}
			if showFile {
				if !v.Show.FilePath {
					file = filepath.Base(file)
				}
				if showPkg || showFunc {
					fileAndLine = fmt.Sprintf(", %s:%d", file, lineNum)
				} else {
					fileAndLine = fmt.Sprintf("%s:%d", file, lineNum)
				}
			}
		}
	}

	if msgFmt != "" { // If there is a message
		if len(args) > 0 {
			msgFmt = fmt.Sprintf(msgFmt, args...)
		}
		if showPkg || showFunc || showFile {
			return fmt.Sprintf("%s%s%s (%s)", pkg, fn, fileAndLine, msgFmt)
		}
		return fmt.Sprintf("(%s)", msgFmt)
	}
	return fmt.Sprintf("%s%s%s", pkg, fn, fileAndLine)
}

func Wrap(e error, v *ErrorDetails, msgFmt string, args ...interface{}) error {
	msg := Format(v, 2, msgFmt, args...)
	return fmt.Errorf("%s%s%w", msg, Sep, e)
}
