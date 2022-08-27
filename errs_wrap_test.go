// Test of errs.Wrap function.
// This test uses HARDCODED LINE NUMBERS for simplicity,
// so it is recommended NOT to insert/remove any lines
// in the middle of this file, but write new test cases appending
// them to the end of TestWrap function.
// If it is required to add/remove lines in the middle,
// re-check all test cases for line number validity.
package errs

import (
	"io"
	"testing"
)

func TestWrap(t *testing.T) {
	ed := NewErrorDetails("test.com/dummy")
	err := io.ErrUnexpectedEOF
	e := Wrap(err, ed, "dummy operation failed")
	msg := e.Error()
	expected := "pkg dummy.TestWrap, errs_wrap_test.go:18 (dummy operation failed): ->\n    <- unexpected EOF"
	expect(t, msg == expected, "default case msg==expected failed, %q != %q", msg, expected)
}
