// Redis serialization protocol specification (RESP)
// Is a  wire protocol that clients implement
// To communicate with the Redis server.

package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Writing RESP
// define constants that represent each type.
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// define a struct to use in the serialization and deserialization process.
// which will hold all the commands and arguments we receive from the client.
type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

// The Reader
type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// readLine reads the line from the buffer.
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

// readInteger reads the integer from the buffer.

func (r *Resp) readInteger() (x int, n int, err error) {
	// read a line and parse it as integer
	line, n, err := r.readLine()
	if err != nil {
		return 0, n, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

// Parsing
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown RESP type: %v\n", string(_type))
		return Value{}, nil
	}
}

// read the Array --Parsing
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// read length of array
	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// foreach line, parse and read the value
	v.array = make([]Value, length)
	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		// add parsed value to array
		v.array[i] = val
	}

	return v, nil
}

// readBulk -- Parsing
func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.typ = "bulk"

	// read length of bulk string
	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, length)

	// read bulk data from reader
	_, err = io.ReadFull(r.reader, bulk)
	if err != nil {
		return v, err
	}

	v.bulk = string(bulk)

	// Read the trailing CRLF
	_, _, err = r.readLine()
	if err != nil {
		return v, err
	}

	return v, nil
}

// Writing RESP (Writing the Value Serializer)
// Marshal method, which will call the specific method for each type based on the Value type.

func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "error":
		return v.marshalError()
	case "null":
		return v.marshalNull()
	default:
		return []byte{}
	}
}

// Simple Strings
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// Bulk String
func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// Array
func (v Value) marshalArray() []byte {
	length := len(v.array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(length)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < length; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}

	return bytes
}

// Null and Error
func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

// Writer --> writing the bytes from Marshal Method to the writer

// Writer struct that takes io.Writer.
type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

// method that takes Value and writes the bytes it gets from the Marshal method to the Writer

func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
