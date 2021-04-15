package util

import (
	"io"
	"sync"
)

// Simultaneous write meanwhile we are reading
func Write(out io.Writer, c chan []byte, lock *sync.WaitGroup) {
	for v := range c {
		_, err := out.Write(v)
		if err != nil {
			panic("ERROR WRITING")
		}
	}

	lock.Done()
}
