// Test of errs.Format function.
// This test uses HARDCODED LINE NUMBERS for simplicity,
// so it is recommended NOT to insert/remove any lines
// in the middle of this file, but write new test cases appending
// them to the end of TestFormat function.
// If it is required to add/remove lines in the middle,
// re-check all test cases for line number validity.
package errs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var expect = require.True

func TestFormat(t *testing.T) {
	ed := NewErrorDetails("test.com/dummy")

	// case with nil ErrorDetails and empty message
	msg := Format(nil, 1, "")
	expected := "[!ErrorDetails==nil] in TestFormat, errs_format_test.go:23"
	expect(t, msg == expected,
		"case with nil ErrorDetails and empty message failed, %q != %q", msg, expected)

	// nil ErrorDetails and constant message
	msg = Format(nil, 1, "test_message")
	expected = "[!ErrorDetails==nil] in TestFormat, errs_format_test.go:29 (test_message)"
	expect(t, msg == expected,
		"case with nil ErrorDetails and constant message failed, %q != %q", msg, expected)

	// nil ErrorDetails and formatted message
	msg = Format(nil, 1, "test_message #%d", 2)
	expected = "[!ErrorDetails==nil] in TestFormat, errs_format_test.go:35 (test_message #2)"
	expect(t, msg == expected,
		"default case with nil ErrorDetails and formatted message failed, %q != %q", msg, expected)

	// case with default ErrorDetails and empty message
	msg = Format(ed, 1, "")
	expected = "pkg dummy.TestFormat, errs_format_test.go:41"
	expect(t, msg == expected, "default case with empty message failed, %q != %q", msg, expected)

	// default ErrorDetails and constant message
	msg = Format(ed, 1, "test_message")
	expected = "pkg dummy.TestFormat, errs_format_test.go:46 (test_message)"
	expect(t, msg == expected, "default case with constant message failed, %q != %q", msg, expected)

	// default ErrorDetails and formatted message
	msg = Format(ed, 1, "test_message #%d", 3)
	expected = "pkg dummy.TestFormat, errs_format_test.go:51 (test_message #3)"
	expect(t, msg == expected, "default case with formatted message failed, %q != %q", msg, expected)

	// full package name
	ed.Show.FullPackageName = true
	msg = Format(ed, 1, "test_message")
	expected = "pkg test.com/dummy.TestFormat, errs_format_test.go:57 (test_message)"
	expect(t, msg == expected, "full package name case failed, %q != %q", msg, expected)

	// full package name + file path
	ed.Show.FilePath = true
	msg = Format(ed, 1, "test_message")
	expectedPrefix := "pkg test.com/dummy.TestFormat, "
	expectedSuffix := "/errs_format_test.go:63 (test_message)"
	expect(t, strings.HasPrefix(msg, expectedPrefix),
		"full package name  + full file name case failed, %q does not have prefix %q", msg, expectedPrefix)
	expect(t, strings.HasSuffix(msg, expectedSuffix),
		"full package name  + full file name case failed, %q does not have suffix %q", msg, expectedSuffix)
	expectedMinLength := len(expectedPrefix) + len(expectedSuffix)
	expect(t, len(msg) >= expectedMinLength,
		"full package name  + full file name case failed, msg %q has length less than expectedMinLength %d", msg, expectedMinLength)

	// hide package name
	ed.Show.Pkg = false
	ed.Show.FilePath = false
	msg = Format(ed, 1, "test_message")
	expected = "TestFormat, errs_format_test.go:77 (test_message)"
	expect(t, msg == expected,
		"hide package name case failed, %q != %q", msg, expected)

	// show only file name
	ed.Show.Func = false
	msg = Format(ed, 1, "test_message")
	expected = "errs_format_test.go:84 (test_message)"
	expect(t, msg == expected,
		"show only file name case failed, %q != %q", msg, expected)

	// hide everything except message
	ed.Show.File = false
	msg = Format(ed, 1, "test_message")
	expected = "(test_message)"
	expect(t, msg == expected,
		"hide everything except message case failed, %q != %q", msg, expected)
}
