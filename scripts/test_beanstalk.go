package main

import (
	"fmt"
	"time"

	"github.com/kr/beanstalk"
)

// Example sends a message
func Example() {
	c, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	if err != nil {
		panic(err)
	}
	c.Put([]byte("hello"), 1, 0, 120*time.Second)
	id, body, err := c.Reserve(5 * time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println("[Example] job", id)
	fmt.Println("[Example]", string(body))
}

// ExampleTubeSetReserve waits for a message and aknowledge
func ExampleTubeSetReserve() {
	c, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	if err != nil {
		panic(err)
	}
	id, body, err := c.Reserve(10 * time.Hour)
	if cerr, ok := err.(beanstalk.ConnError); ok && cerr.Err == beanstalk.ErrTimeout {
		fmt.Println("timed out")
		return
	} else if err != nil {
		panic(err)
	}
	fmt.Println("[ExampleTubeSet_Reserve] job", id)
	fmt.Println("[ExampleTubeSet_Reserve]", string(body))
}

func main() {
	go ExampleTubeSetReserve()
	Example()
}
