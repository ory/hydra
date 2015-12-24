// Package redis is a simple client for connecting and interacting with a single
// redis instance.
//
// THE FUNCTIONALITY PROVIDED IN THIS PACKAGE IS NOT THREAD-SAFE. To use a
// single redis instance amongst multiple go-routines, check out the pool
// subpackage (http://godoc.org/github.com/mediocregopher/radix.v2/pool)
//
// To import inside your package do:
//
//	import "github.com/mediocregopher/radix.v2/redis"
//
// Connecting
//
// Use either Dial or DialTimeout:
//
//	client, err := redis.Dial("tcp", "localhost:6379")
//	if err != nil {
//		// handle err
//	}
//
// Make sure to call Close on the client if you want to clean it up before the
// end of the program.
//
// Cmd and Resp
//
// The Cmd method returns a Resp, which has methods for converting to various
// types. Each of these methods returns an error which can either be a
// connection error (e.g. timeout), an application error (e.g. key is wrong
// type), or a conversion error (e.g. cannot convert to integer). You can also
// directly check the error using the Err field:
//
//	foo, err := client.Cmd("GET", "foo").Str()
//	if err != nil {
//		// handle err
//	}
//
//	// Checking Err field directly
//
//	err = client.Cmd("SET", "foo", "bar", "EX", 3600).Err
//	if err != nil {
//		// handle err
//	}
//
// Array Replies
//
// The elements to Array replies can be accessed as strings using List or
// ListBytes, or you can use the Array method for more low level access:
//
//	r := client.Cmd("MGET", "foo", "bar", "baz")
//	if r.Err != nil {
//		// handle error
//	}
//
//	// This:
//	l, _ := r.List()
//	for _, elemStr := range l {
//		fmt.Println(elemStr)
//	}
//
//	// is equivalent to this:
//	elems, err := r.Array()
//	for i := range elems {
//		elemStr, _ := elems[i].Str()
//		fmt.Println(elemStr)
//	}
//
// Pipelining
//
// Pipelining is when the client sends a bunch of commands to the server at
// once, and only once all the commands have been sent does it start reading the
// replies off the socket. This is supported using the PipeAppend and PipeResp
// methods. PipeAppend will simply append the command to a buffer without
// sending it, the first time PipeResp is called it will send all the commands
// in the buffer and return the Resp for the first command that was sent.
// Subsequent calls to PipeResp return Resps for subsequent commands:
//
//	client.PipeAppend("GET", "foo")
//	client.PipeAppend("SET", "bar", "foo")
//	client.PipeAppend("DEL", "baz")
//
//	// Read GET foo reply
//	foo, err := client.PipeResp().Str()
//	if err != nil {
//		// handle err
//	}
//
//	// Read SET bar foo reply
//	if err := client.PipeResp().Err; err != nil {
//		// handle err
//	}
//
//	// Read DEL baz reply
//	if err := client.PipeResp().Err; err != nil {
//		// handle err
//	}
//
// Flattening
//
// Radix will automatically flatten passed in maps and slices into the argument
// list. For example, the following are all equivalent:
//
//	client.Cmd("HMSET", "myhash", "key1", "val1", "key2", "val2")
//	client.Cmd("HMSET", "myhash", []string{"key1", "val1", "key2", "val2"})
//	client.Cmd("HMSET", "myhash", map[string]string{
//		"key1": "val1",
//		"key2": "val2",
//	})
//	client.Cmd("HMSET", "myhash", [][]string{
//		[]string{"key1", "val1"},
//		[]string{"key2", "val2"},
//	})
//
// Radix is not picky about the types inside or outside the maps/slices, if they
// don't match a subset of primitive types it will fall back to reflection to
// figure out what they are and encode them.
package redis
