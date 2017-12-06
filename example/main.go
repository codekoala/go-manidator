package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/codekoala/go-manidator"
)

func main() {
	// handle keyboard interrupts
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	// be a good citizen
	ctx, cancel := context.WithCancel(context.Background())

	// make some things
	b1 := newDator("fred")
	b2 := newDator("wilma")
	b3 := newDator("barney")

	// prepare to output things from the things
	m := manidator.New(b1, b2, b3)

	fmt.Println("starting")
	m.Begin(ctx)

	// wait for things to happen
	select {
	case <-stop:
		fmt.Println("stopping")
		cancel()
	case <-m.Done():
		fmt.Println("all dators finished")
	}

	// stop the things for good measure
	m.Stop()
	fmt.Println("finished")
}

// newDator creates a new Dator and immediately begins reading lines from the included lipsum.txt file.
func newDator(name string) manidator.Dator {
	d := manidator.NewDator(name)

	go randReader(d)

	return d
}

// randReader reads lines from the included lipsum.txt file, pausing for a random number of milliseconds between each
// line.
func randReader(dest io.WriteCloser) {
	defer dest.Close()

	file, err := os.Open("lipsum.txt")
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		dest.Write(sc.Bytes())
		dest.Write([]byte("\n"))

		r := rand.Intn(200)
		time.Sleep(time.Duration(r+50) * time.Millisecond)
	}

	if err = sc.Err(); err != nil {
		fmt.Printf("SCAN ERROR: %s\n", err)
	}
}
