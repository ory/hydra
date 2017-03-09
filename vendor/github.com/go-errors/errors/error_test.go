package errors

import (
	"bytes"
	"fmt"
	"io"
	"runtime/debug"
	"testing"
)

func TestStackFormatMatches(t *testing.T) {

	defer func() {
		err := recover()
		if err != 'a' {
			t.Fatal(err)
		}

		bs := [][]byte{Errorf("hi").Stack(), debug.Stack()}

		// Ignore the first line (as it contains the PC of the .Stack() call)
		bs[0] = bytes.SplitN(bs[0], []byte("\n"), 2)[1]
		bs[1] = bytes.SplitN(bs[1], []byte("\n"), 2)[1]

		if bytes.Compare(bs[0], bs[1]) != 0 {
			t.Errorf("Stack didn't match")
			t.Errorf("%s", bs[0])
			t.Errorf("%s", bs[1])
		}
	}()

	a()
}

func TestSkipWorks(t *testing.T) {

	defer func() {
		err := recover()
		if err != 'a' {
			t.Fatal(err)
		}

		bs := [][]byte{Wrap("hi", 2).Stack(), debug.Stack()}

		// should skip four lines of debug.Stack()
		bs[1] = bytes.SplitN(bs[1], []byte("\n"), 5)[4]

		if bytes.Compare(bs[0], bs[1]) != 0 {
			t.Errorf("Stack didn't match")
			t.Errorf("%s", bs[0])
			t.Errorf("%s", bs[1])
		}
	}()

	a()
}

func TestNew(t *testing.T) {

	err := New("foo")

	if err.Error() != "foo" {
		t.Errorf("Wrong message")
	}

	err = New(fmt.Errorf("foo"))

	if err.Error() != "foo" {
		t.Errorf("Wrong message")
	}

	bs := [][]byte{New("foo").Stack(), debug.Stack()}

	// Ignore the first line (as it contains the PC of the .Stack() call)
	bs[0] = bytes.SplitN(bs[0], []byte("\n"), 2)[1]
	bs[1] = bytes.SplitN(bs[1], []byte("\n"), 2)[1]

	if bytes.Compare(bs[0], bs[1]) != 0 {
		t.Errorf("Stack didn't match")
		t.Errorf("%s", bs[0])
		t.Errorf("%s", bs[1])
	}

	if err.ErrorStack() != err.TypeName()+" "+err.Error()+"\n"+string(err.Stack()) {
		t.Errorf("ErrorStack is in the wrong format")
	}
}

func TestIs(t *testing.T) {

	if Is(nil, io.EOF) {
		t.Errorf("nil is an error")
	}

	if !Is(io.EOF, io.EOF) {
		t.Errorf("io.EOF is not io.EOF")
	}

	if !Is(io.EOF, New(io.EOF)) {
		t.Errorf("io.EOF is not New(io.EOF)")
	}

	if !Is(New(io.EOF), New(io.EOF)) {
		t.Errorf("New(io.EOF) is not New(io.EOF)")
	}

	if Is(io.EOF, fmt.Errorf("io.EOF")) {
		t.Errorf("io.EOF is fmt.Errorf")
	}

}

func TestWrapError(t *testing.T) {

	e := func() error {
		return Wrap("hi", 1)
	}()

	if e.Error() != "hi" {
		t.Errorf("Constructor with a string failed")
	}

	if Wrap(fmt.Errorf("yo"), 0).Error() != "yo" {
		t.Errorf("Constructor with an error failed")
	}

	if Wrap(e, 0) != e {
		t.Errorf("Constructor with an Error failed")
	}

	if Wrap(nil, 0).Error() != "<nil>" {
		t.Errorf("Constructor with nil failed")
	}
}

func TestWrapPrefixError(t *testing.T) {

	e := func() error {
		return WrapPrefix("hi", "prefix", 1)
	}()

	fmt.Println(e.Error())
	if e.Error() != "prefix: hi" {
		t.Errorf("Constructor with a string failed")
	}

	if WrapPrefix(fmt.Errorf("yo"), "prefix", 0).Error() != "prefix: yo" {
		t.Errorf("Constructor with an error failed")
	}

	prefixed := WrapPrefix(e, "prefix", 0)
	original := e.(*Error)

	if prefixed.Err != original.Err || &prefixed.stack != &original.stack || &prefixed.frames != &original.frames || prefixed.Error() != "prefix: prefix: hi" {
		t.Errorf("Constructor with an Error failed")
	}

	if WrapPrefix(nil, "prefix", 0).Error() != "prefix: <nil>" {
		t.Errorf("Constructor with nil failed")
	}
}

func ExampleErrorf(x int) (int, error) {
	if x%2 == 1 {
		return 0, Errorf("can only halve even numbers, got %d", x)
	}
	return x / 2, nil
}

func ExampleWrapError() (error, error) {
	// Wrap io.EOF with the current stack-trace and return it
	return nil, Wrap(io.EOF, 0)
}

func ExampleWrapError_skip() {
	defer func() {
		if err := recover(); err != nil {
			// skip 1 frame (the deferred function) and then return the wrapped err
			err = Wrap(err, 1)
		}
	}()
}

func ExampleIs(reader io.Reader, buff []byte) {
	_, err := reader.Read(buff)
	if Is(err, io.EOF) {
		return
	}
}

func ExampleNew(UnexpectedEOF error) error {
	// calling New attaches the current stacktrace to the existing UnexpectedEOF error
	return New(UnexpectedEOF)
}

func ExampleWrap() error {

	if err := recover(); err != nil {
		return Wrap(err, 1)
	}

	return a()
}

func ExampleError_Error(err error) {
	fmt.Println(err.Error())
}

func ExampleError_ErrorStack(err error) {
	fmt.Println(err.(*Error).ErrorStack())
}

func ExampleError_Stack(err *Error) {
	fmt.Println(err.Stack())
}

func ExampleError_TypeName(err *Error) {
	fmt.Println(err.TypeName(), err.Error())
}

func ExampleError_StackFrames(err *Error) {
	for _, frame := range err.StackFrames() {
		fmt.Println(frame.File, frame.LineNumber, frame.Package, frame.Name)
	}
}

func a() error {
	b(5)
	return nil
}

func b(i int) {
	c()
}

func c() {
	panic('a')
}
