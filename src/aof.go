package main

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"
)

// Data persistence
// Write the AOF(Append only file) to disk

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

// NewAof method which is used in main.go upon server start.
func NewAof(path string) (*Aof, error) {

	// Open the aof file
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	// create buflio reader to read from the file
	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// Start a goroutine to sync AOF to disk every 1 second
	go func() {
		for {
			aof.mu.Lock()
			aof.file.Sync()
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

// Close method to close the AOF file
func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

// Write method to write commands to the aof file
func (aof *Aof) Write(value Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal()) // to write the command to the file in the same RESP format that we receive
	if err != nil {
		return err
	}

	return nil
}

// Read method to read commands from the aof file
func (aof *Aof) Read(callback func(value Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	resp := NewResp(aof.file)

	for {
		value, err := resp.Read()
		if err == nil {
			callback(value)
		}
		if err == io.EOF {
			break
		}
		return err
	}

	return nil
}
