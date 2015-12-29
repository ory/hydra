package redis

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"strconv"
	"strings"
)

var (
	delim    = []byte{'\r', '\n'}
	delimEnd = delim[len(delim)-1]
)

// RespType is a field on every Resp which indicates the type of the data it
// contains
type RespType int

// Different RespTypes. You can check if a message is of one or more types using
// the IsType method on Resp
const (
	SimpleStr RespType = 1 << iota
	BulkStr
	IOErr  // An error which prevented reading/writing, e.g. connection close
	AppErr // An error returned by redis, e.g. WRONGTYPE
	Int
	Array
	Nil

	// Str combines both SimpleStr and BulkStr, which are considered strings to
	// the Str() method.  This is what you want to give to IsType when
	// determining if a response is a string
	Str = SimpleStr | BulkStr

	// Err combines both IOErr and AppErr, which both indicate that the Err
	// field on their Resp is filled. To determine if a Resp is an error you'll
	// most often want to simply check if the Err field on it is nil
	Err = IOErr | AppErr
)

var (
	simpleStrPrefix = []byte{'+'}
	errPrefix       = []byte{'-'}
	intPrefix       = []byte{':'}
	bulkStrPrefix   = []byte{'$'}
	arrayPrefix     = []byte{'*'}
	nilFormatted    = []byte("$-1\r\n")
)

// Parse errors
var (
	errBadType  = errors.New("wrong type")
	errParse    = errors.New("parse error")
	errNotStr   = errors.New("could not convert to string")
	errNotInt   = errors.New("could not convert to int")
	errNotArray = errors.New("could not convert to array")
)

// Resp represents a single response or message being sent to/from a redis
// server. Each Resp has a type (see RespType and IsType) and a value. Values
// can be retrieved using any of the casting methods on this type (e.g. Str)
type Resp struct {
	typ RespType
	val interface{}

	// Err indicates that this Resp signals some kind of error, either on the
	// connection level or the application level. Use IsType if you need to
	// determine which, otherwise you can simply check if this is nil
	Err error
}

// NewResp takes the given value and interprets it into a resp encoded byte
// stream
func NewResp(v interface{}) *Resp {
	r := format(v, false)
	return &r
}

// NewRespSimple is like NewResp except it encodes its string as a resp
// SimpleStr type, whereas NewResp will encode all strings as BulkStr
func NewRespSimple(s string) *Resp {
	return &Resp{
		typ: SimpleStr,
		val: []byte(s),
	}
}

// NewRespFlattenedStrings is like NewResp except it looks through the given
// value and converts any types (except slices/maps) into strings, and flatten any
// embedded slices/maps into a single slice. This is useful because commands to
// a redis server must be given as an array of bulk strings. If the argument
// isn't already in a slice/map it will be wrapped so that it is written as a
// Array of size one
func NewRespFlattenedStrings(v interface{}) *Resp {
	fv := flatten(v)
	r := format(fv, true)
	return &r
}

// newRespIOErr is a convenience method for making Resps to wrap io errors
func newRespIOErr(err error) *Resp {
	r := NewResp(err)
	r.typ = IOErr
	return r
}

// RespReader is a wrapper around an io.Reader which will read Resp messages off
// of the io.Reader
type RespReader struct {
	r *bufio.Reader
}

// NewRespReader creates and returns a new RespReader which will read from the
// given io.Reader. Once passed in the io.Reader shouldn't be read from by any
// other processes
func NewRespReader(r io.Reader) *RespReader {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	return &RespReader{br}
}

// ReadResp attempts to read a message object from the given io.Reader, parse
// it, and return a Resp representing it
func (rr *RespReader) Read() *Resp {
	res, err := bufioReadResp(rr.r)
	if err != nil {
		res = Resp{typ: IOErr, val: err, Err: err}
	}
	return &res
}

func bufioReadResp(r *bufio.Reader) (Resp, error) {
	b, err := r.Peek(1)
	if err != nil {
		return Resp{}, err
	}
	switch b[0] {
	case simpleStrPrefix[0]:
		return readSimpleStr(r)
	case errPrefix[0]:
		return readError(r)
	case intPrefix[0]:
		return readInt(r)
	case bulkStrPrefix[0]:
		return readBulkStr(r)
	case arrayPrefix[0]:
		return readArray(r)
	default:
		return Resp{}, errBadType
	}
}

func readSimpleStr(r *bufio.Reader) (Resp, error) {
	b, err := r.ReadBytes(delimEnd)
	if err != nil {
		return Resp{}, err
	}
	return Resp{typ: SimpleStr, val: b[1 : len(b)-2]}, nil
}

func readError(r *bufio.Reader) (Resp, error) {
	b, err := r.ReadBytes(delimEnd)
	if err != nil {
		return Resp{}, err
	}
	err = errors.New(string(b[1 : len(b)-2]))
	return Resp{typ: AppErr, val: err, Err: err}, nil
}

func readInt(r *bufio.Reader) (Resp, error) {
	b, err := r.ReadBytes(delimEnd)
	if err != nil {
		return Resp{}, err
	}
	i, err := strconv.ParseInt(string(b[1:len(b)-2]), 10, 64)
	if err != nil {
		return Resp{}, errParse
	}
	return Resp{typ: Int, val: i}, nil
}

func readBulkStr(r *bufio.Reader) (Resp, error) {
	b, err := r.ReadBytes(delimEnd)
	if err != nil {
		return Resp{}, err
	}
	size, err := strconv.ParseInt(string(b[1:len(b)-2]), 10, 64)
	if err != nil {
		return Resp{}, errParse
	}
	if size < 0 {
		return Resp{typ: Nil}, nil
	}
	total := make([]byte, size)
	b2 := total
	var n int
	for len(b2) > 0 {
		n, err = r.Read(b2)
		if err != nil {
			return Resp{}, err
		}
		b2 = b2[n:]
	}

	// There's a hanging \r\n there, gotta read past it
	trail := make([]byte, 2)
	for i := 0; i < 2; i++ {
		c, err := r.ReadByte()
		if err != nil {
			return Resp{}, err
		}
		trail[i] = c
	}

	return Resp{typ: BulkStr, val: total}, nil
}

func readArray(r *bufio.Reader) (Resp, error) {
	b, err := r.ReadBytes(delimEnd)
	if err != nil {
		return Resp{}, err
	}
	size, err := strconv.ParseInt(string(b[1:len(b)-2]), 10, 64)
	if err != nil {
		return Resp{}, errParse
	}
	if size < 0 {
		return Resp{typ: Nil}, nil
	}

	arr := make([]Resp, size)
	for i := range arr {
		m, err := bufioReadResp(r)
		if err != nil {
			return Resp{}, err
		}
		arr[i] = m
	}
	return Resp{typ: Array, val: arr}, nil
}

// IsType returns whether or or not the reply is of a given type
//
//	isStr := r.IsType(redis.Str)
//
// Multiple types can be checked at the same time by or'ing the desired types
//
//	isStrOrInt := r.IsType(redis.Str | redis.Int)
//
func (r *Resp) IsType(t RespType) bool {
	return r.typ&t > 0
}

// WriteTo writes the resp encoded form of the Resp to the given writer,
// implementing the WriterTo interface
func (r *Resp) WriteTo(w io.Writer) (int64, error) {

	// SimpleStr is a special case, writeTo always writes strings as BulkStrs,
	// so we just manually do SimpleStr here
	if r.typ == SimpleStr {
		s := r.val.([]byte)
		b := append(make([]byte, 0, len(s)+3), simpleStrPrefix...)
		b = append(b, s...)
		b = append(b, delim...)
		written, err := w.Write(b)
		return int64(written), err
	}

	return writeTo(w, nil, r.val, false, false)
}

// Bytes returns a byte slice representing the value of the Resp. Only valid for
// a Resp of type Str. If r.Err != nil that will be returned.
func (r *Resp) Bytes() ([]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	} else if !r.IsType(Str) {
		return nil, errBadType
	}

	if b, ok := r.val.([]byte); ok {
		return b, nil
	}
	return nil, errNotStr
}

// Str is a wrapper around Bytes which returns the result as a string instead of
// a byte slice
func (r *Resp) Str() (string, error) {
	b, err := r.Bytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Int returns an int representing the value of the Resp. For a Resp of type Int
// the integer value will be returned directly. For a Resp of type Str the
// string will attempt to be parsed as a base-10 integer, returning the parsing
// error if any. If r.Err != nil that will be returned
func (r *Resp) Int() (int, error) {
	i, err := r.Int64()
	return int(i), err
}

// Int64 is like Int, but returns int64 instead of Int
func (r *Resp) Int64() (int64, error) {
	if r.Err != nil {
		return 0, r.Err
	}
	if i, ok := r.val.(int64); ok {
		return i, nil
	}
	if s, err := r.Str(); err == nil {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, errNotInt
}

// Float64 returns a float64 representing the value of the Resp. Only valud for
// a Resp of type Str which represents an actual float. If r.Err != nil that
// will be returned
func (r *Resp) Float64() (float64, error) {
	if r.Err != nil {
		return 0, r.Err
	}
	if b, ok := r.val.([]byte); ok {
		f, err := strconv.ParseFloat(string(b), 64)
		if err != nil {
			return 0, err
		}
		return f, nil
	}
	return 0, errNotStr
}

func (r *Resp) betterArray() ([]Resp, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	if a, ok := r.val.([]Resp); ok {
		return a, nil
	}
	return nil, errNotArray
}

// Array returns the Resp slice encompassed by this Resp. Only valid for a Resp
// of type Array. If r.Err != nil that will be returned
func (r *Resp) Array() ([]*Resp, error) {
	a, err := r.betterArray()
	if err != nil {
		return nil, err
	}
	abad := make([]*Resp, len(a))
	for i := range a {
		abad[i] = &a[i]
	}
	return abad, nil
}

// List is a wrapper around Array which returns the result as a list of strings,
// calling Str() on each Resp which Array returns. Any errors encountered are
// immediately returned. Any Nil replies are interpreted as empty strings
func (r *Resp) List() ([]string, error) {
	m, err := r.betterArray()
	if err != nil {
		return nil, err
	}
	l := make([]string, len(m))
	for i := range m {
		if m[i].IsType(Nil) {
			l[i] = ""
			continue
		}
		s, err := m[i].Str()
		if err != nil {
			return nil, err
		}
		l[i] = s
	}
	return l, nil
}

// ListBytes is a wrapper around Array which returns the result as a list of
// byte slices, calling Bytes() on each Resp which Array returns. Any errors
// encountered are immediately returned. Any Nil replies are interpreted as nil
func (r *Resp) ListBytes() ([][]byte, error) {
	m, err := r.betterArray()
	if err != nil {
		return nil, err
	}
	l := make([][]byte, len(m))
	for i := range m {
		if m[i].IsType(Nil) {
			l[i] = nil
			continue
		}
		b, err := m[i].Bytes()
		if err != nil {
			return nil, err
		}
		l[i] = b
	}
	return l, nil
}

// Map is a wrapper around Array which returns the result as a map of strings,
// calling Str() on alternating key/values for the map. All value fields of type
// Nil will be treated as empty strings, keys must all be of type Str
func (r *Resp) Map() (map[string]string, error) {
	l, err := r.betterArray()
	if err != nil {
		return nil, err
	}
	if len(l)%2 != 0 {
		return nil, errors.New("reply has odd number of elements")
	}

	m := map[string]string{}
	for {
		if len(l) == 0 {
			return m, nil
		}
		k, v := l[0], l[1]
		l = l[2:]

		ks, err := k.Str()
		if err != nil {
			return nil, err
		}

		var vs string
		if v.IsType(Nil) {
			vs = ""
		} else if vs, err = v.Str(); err != nil {
			return nil, err
		}
		m[ks] = vs
	}
}

// String returns a string representation of the Resp. This method is for
// debugging, use Str() for reading a Str reply
func (r *Resp) String() string {
	var inner string
	switch r.typ {
	case AppErr:
		inner = fmt.Sprintf("AppErr %s", r.Err)
	case IOErr:
		inner = fmt.Sprintf("IOErr %s", r.Err)
	case BulkStr, SimpleStr:
		inner = fmt.Sprintf("Str %q", string(r.val.([]byte)))
	case Int:
		inner = fmt.Sprintf("Int %d", r.val.(int64))
	case Nil:
		inner = fmt.Sprintf("Nil")
	case Array:
		kids := r.val.([]Resp)
		kidsStr := make([]string, len(kids))
		for i := range kids {
			kidsStr[i] = kids[i].String()
		}
		inner = strings.Join(kidsStr, " ")
	default:
		inner = "UNKNOWN"
	}
	return fmt.Sprintf("Resp(%s)", inner)
}

var typeOfBytes = reflect.TypeOf([]byte(nil))

func flattenedLength(mm ...interface{}) int {

	total := 0

	for _, m := range mm {
		switch m.(type) {
		case []byte, string, bool, nil, int, int8, int16, int32, int64, uint,
			uint8, uint16, uint32, uint64, float32, float64, error:
			total++

		case Resp:
			total += flattenedLength(m.(Resp).val)
		case *Resp:
			total += flattenedLength(m.(*Resp).val)

		case []interface{}:
			total += flattenedLength(m.([]interface{})...)

		default:
			t := reflect.TypeOf(m)

			switch t.Kind() {
			case reflect.Slice:
				rm := reflect.ValueOf(m)
				l := rm.Len()
				for i := 0; i < l; i++ {
					total += flattenedLength(rm.Index(i).Interface())
				}

			case reflect.Map:
				rm := reflect.ValueOf(m)
				keys := rm.MapKeys()
				for _, k := range keys {
					kv := k.Interface()
					vv := rm.MapIndex(k).Interface()
					total += flattenedLength(kv)
					total += flattenedLength(vv)
				}

			default:
				total++
			}
		}
	}

	return total
}

func flatten(m interface{}) []interface{} {
	t := reflect.TypeOf(m)

	// If it's a byte-slice we don't want to flatten
	if t == typeOfBytes {
		return []interface{}{m}
	}

	switch t.Kind() {
	case reflect.Slice:
		rm := reflect.ValueOf(m)
		l := rm.Len()
		ret := make([]interface{}, 0, l)
		for i := 0; i < l; i++ {
			ret = append(ret, flatten(rm.Index(i).Interface())...)
		}
		return ret

	case reflect.Map:
		rm := reflect.ValueOf(m)
		l := rm.Len() * 2
		keys := rm.MapKeys()
		ret := make([]interface{}, 0, l)
		for _, k := range keys {
			kv := k.Interface()
			vv := rm.MapIndex(k).Interface()
			ret = append(ret, flatten(kv)...)
			ret = append(ret, flatten(vv)...)
		}
		return ret

	default:
		return []interface{}{m}
	}
}

func anyIntToInt64(m interface{}) int64 {
	switch mt := m.(type) {
	case int:
		return int64(mt)
	case int8:
		return int64(mt)
	case int16:
		return int64(mt)
	case int32:
		return int64(mt)
	case int64:
		return mt
	case uint:
		return int64(mt)
	case uint8:
		return int64(mt)
	case uint16:
		return int64(mt)
	case uint32:
		return int64(mt)
	case uint64:
		return int64(mt)
	}
	panic(fmt.Sprintf("anyIntToInt64 got bad arg: %#v", m))
}

func writeBytesHelper(
	w io.Writer, b []byte, lastWritten int64, lastErr error,
) (
	int64, error,
) {
	if lastErr != nil {
		return lastWritten, lastErr
	}
	i, err := w.Write(b)
	return int64(i) + lastWritten, err
}

func writeArrayHeader(w io.Writer, buf []byte, l int64) (int64, error) {
	buf = strconv.AppendInt(buf, l, 10)
	var err error
	var written int64
	written, err = writeBytesHelper(w, arrayPrefix, written, err)
	written, err = writeBytesHelper(w, buf, written, err)
	written, err = writeBytesHelper(w, delim, written, err)
	return written, err
}

// Given a preallocated byte buffer and a string, this will copy the string's
// contents into buf starting at index 0, and returns two slices from buf: The
// first is a slice of the string data, the second is a slice of the "rest" of
// buf following the first slice
func stringSlicer(buf []byte, s string) ([]byte, []byte) {
	sbuf := append(buf[:0], s...)
	return sbuf, sbuf[len(sbuf):]
}

// takes in something, m, and encodes it and writes it to w. buf is used as a
// pre-alloated byte buffer for encoding integers (expected to have a length of
// 0), so we don't have to re-allocate a new one every time we convert an
// integer to a string.  forceString means all types will be converted to
// strings, noArrayHeader means don't write out the headers to any arrays, just
// inline all the elements in the array
func writeTo(
	w io.Writer, buf []byte, m interface{}, forceString, noArrayHeader bool,
) (
	int64, error,
) {
	switch mt := m.(type) {
	case []byte:
		return writeStr(w, buf, mt)
	case string:
		sbuf, buf := stringSlicer(buf, mt)
		return writeStr(w, buf, sbuf)
	case bool:
		buf = buf[:0]
		if mt {
			buf = append(buf, '1')
		} else {
			buf = append(buf, '0')
		}
		return writeStr(w, buf[1:], buf[:1])
	case nil:
		if forceString {
			return writeStr(w, buf, nil)
		}
		return writeNil(w)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		i := anyIntToInt64(mt)
		return writeInt(w, buf, i, forceString)
	case float32:
		return writeFloat(w, buf, float64(mt), 32)
	case float64:
		return writeFloat(w, buf, mt, 64)
	case error:
		return writeErr(w, buf, mt, forceString)

	// We duplicate the below code here a bit, since this is the common case and
	// it'd be better to not get the reflect package involved here
	case []interface{}:
		l := len(mt)
		var totalWritten int64

		if !noArrayHeader {
			written, err := writeArrayHeader(w, buf, int64(l))
			totalWritten += written
			if err != nil {
				return totalWritten, err
			}
		}
		for i := 0; i < l; i++ {
			written, err := writeTo(w, buf, mt[i], forceString, noArrayHeader)
			totalWritten += written
			if err != nil {
				return totalWritten, err
			}
		}
		return totalWritten, nil

	case *Resp:
		return writeTo(w, buf, mt.val, forceString, noArrayHeader)

	case Resp:
		return writeTo(w, buf, mt.val, forceString, noArrayHeader)

	default:
		// Fallback to reflect-based.
		switch reflect.TypeOf(m).Kind() {
		case reflect.Slice:
			rm := reflect.ValueOf(mt)
			l := rm.Len()
			var totalWritten, written int64
			var err error

			if !noArrayHeader {
				written, err = writeArrayHeader(w, buf, int64(l))
				totalWritten += written
				if err != nil {
					return totalWritten, err
				}
			}
			for i := 0; i < l; i++ {
				vv := rm.Index(i).Interface()
				written, err = writeTo(w, buf, vv, forceString, noArrayHeader)
				totalWritten += written
				if err != nil {
					return totalWritten, err
				}
			}
			return totalWritten, nil

		case reflect.Map:
			rm := reflect.ValueOf(mt)
			l := rm.Len() * 2
			var totalWritten, written int64
			var err error

			if !noArrayHeader {
				written, err = writeArrayHeader(w, buf, int64(l))
				totalWritten += written
				if err != nil {
					return totalWritten, err
				}
			}
			keys := rm.MapKeys()
			for _, k := range keys {
				kv := k.Interface()
				written, err = writeTo(w, buf, kv, forceString, noArrayHeader)
				totalWritten += written
				if err != nil {
					return totalWritten, err
				}

				vv := rm.MapIndex(k).Interface()
				written, err = writeTo(w, buf, vv, forceString, noArrayHeader)
				if err != nil {
					return totalWritten, err
				}
			}
			return totalWritten, nil

		default:
			return writeStr(w, buf, []byte(fmt.Sprint(m)))
		}
	}
}

func writeStr(w io.Writer, buf, b []byte) (int64, error) {
	var err error
	var written int64
	buf = strconv.AppendInt(buf[:0], int64(len(b)), 10)

	written, err = writeBytesHelper(w, bulkStrPrefix, written, err)
	written, err = writeBytesHelper(w, buf, written, err)
	written, err = writeBytesHelper(w, delim, written, err)
	written, err = writeBytesHelper(w, b, written, err)
	written, err = writeBytesHelper(w, delim, written, err)
	return written, err
}

func writeErr(
	w io.Writer, buf []byte, ierr error, forceString bool,
) (
	int64, error,
) {
	ierrStr := []byte(ierr.Error())
	if forceString {
		return writeStr(w, buf, ierrStr)
	}
	var err error
	var written int64
	written, err = writeBytesHelper(w, errPrefix, written, err)
	written, err = writeBytesHelper(w, []byte(ierr.Error()), written, err)
	written, err = writeBytesHelper(w, delim, written, err)
	return written, err
}

func writeInt(
	w io.Writer, buf []byte, i int64, forceString bool,
) (
	int64, error,
) {
	buf = strconv.AppendInt(buf[:0], i, 10)
	if forceString {
		return writeStr(w, buf[len(buf):], buf)
	}

	var err error
	var written int64
	written, err = writeBytesHelper(w, intPrefix, written, err)
	written, err = writeBytesHelper(w, buf, written, err)
	written, err = writeBytesHelper(w, delim, written, err)
	return written, err
}

func writeFloat(w io.Writer, buf []byte, f float64, bits int) (int64, error) {
	buf = strconv.AppendFloat(buf[:0], f, 'f', -1, bits)
	return writeStr(w, buf[len(buf):], buf)
}

func writeNil(w io.Writer) (int64, error) {
	written, err := w.Write(nilFormatted)
	return int64(written), err
}

// IsTimeout is a helper function for determining if an IOErr Resp was caused by
// a network timeout
func IsTimeout(r *Resp) bool {
	if r.IsType(IOErr) {
		t, ok := r.Err.(*net.OpError)
		return ok && t.Timeout()
	}
	return false
}

// format takes any data structure and attempts to turn it into a Resp or
// multiple embedded Resps in the form of an Array. This is only used for
// NewResp and NewRespFlattenedStrings
func format(m interface{}, forceString bool) Resp {
	switch mt := m.(type) {
	case []byte:
		return Resp{typ: BulkStr, val: mt}
	case string:
		return Resp{typ: BulkStr, val: []byte(mt)}
	case bool:
		if mt {
			return Resp{typ: BulkStr, val: []byte{'1'}}
		}
		return Resp{typ: BulkStr, val: []byte{'0'}}
	case nil:
		if forceString {
			return Resp{typ: BulkStr, val: []byte{'0'}}
		}
		return Resp{typ: Nil}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		i := anyIntToInt64(mt)
		if forceString {
			istr := strconv.FormatInt(i, 10)
			return Resp{typ: BulkStr, val: []byte(istr)}
		}
		return Resp{typ: Int, val: i}
	case float32:
		ft := strconv.FormatFloat(float64(mt), 'f', -1, 32)
		return Resp{typ: BulkStr, val: []byte(ft)}
	case float64:
		ft := strconv.FormatFloat(mt, 'f', -1, 64)
		return Resp{typ: BulkStr, val: []byte(ft)}
	case error:
		if forceString {
			return Resp{typ: BulkStr, val: []byte(mt.Error())}
		}
		return Resp{typ: AppErr, val: mt, Err: mt}

	// We duplicate the below code here a bit, since this is the common case and
	// it'd be better to not get the reflect package involved here
	case []interface{}:
		l := len(mt)
		rl := make([]Resp, 0, l)
		for i := 0; i < l; i++ {
			r := format(mt[i], forceString)
			rl = append(rl, r)
		}
		return Resp{typ: Array, val: rl}

	case *Resp:
		return *mt

	case Resp:
		return mt

	default:
		// Fallback to reflect-based.
		switch reflect.TypeOf(m).Kind() {
		case reflect.Slice:
			rm := reflect.ValueOf(mt)
			l := rm.Len()
			rl := make([]Resp, 0, l)
			for i := 0; i < l; i++ {
				vv := rm.Index(i).Interface()
				r := format(vv, forceString)
				rl = append(rl, r)
			}
			return Resp{typ: Array, val: rl}

		case reflect.Map:
			rm := reflect.ValueOf(mt)
			l := rm.Len() * 2
			rl := make([]Resp, 0, l)
			keys := rm.MapKeys()
			for _, k := range keys {
				kv := k.Interface()
				vv := rm.MapIndex(k).Interface()

				kr := format(kv, forceString)
				rl = append(rl, kr)

				vr := format(vv, forceString)
				rl = append(rl, vr)
			}
			return Resp{typ: Array, val: rl}

		default:
			return Resp{typ: BulkStr, val: []byte(fmt.Sprint(m))}
		}
	}
}
