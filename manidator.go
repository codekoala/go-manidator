package manidator

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

const ESC = 27

var (
	// DefaultInterval is the default amount of time to wait between each update to the screen.
	DefaultInterval = 50 * time.Millisecond

	// DefaultTerminalWidth is the default size of the output window when it cannot be determined automatically.
	DefaultTerminalWidth = 80

	// eraseLine is a series of control codes that move the cursor up one line and erase that line.
	eraseLine = fmt.Sprintf("%c[%dA%c[2K\r", ESC, 1, ESC)
)

type Manidator struct {
	// Interval is the amount of time to wait between each update to the screen. By default it is 50ms.
	Interval time.Duration

	// dators is a list of Dators which are expected to provide output to the user.
	dators []Dator

	// nameWidth is the number of columns that will be consumed by the name of each Dator. This is determined
	// automatically as soon as Manidator.Begin() is called.
	nameWidth int

	// fmtString is the format string used to print each Dator's last line of output, prefixed by each respective
	// Dator's name.
	fmtString string

	// done is a channel which, when closed, indicates that all Dators have finished providing output, and allows the
	// Manidator to stop itself.
	done chan struct{}

	// cancel is a channel which, when closed, indicates that the Manidator owner has requested that no more updates be
	// provided.
	cancel chan struct{}

	// wg helps keep track of goroutines.
	wg sync.WaitGroup
}

// New creates a new Manidator, optionally adding one or more Dators immediately.
func New(dators ...Dator) *Manidator {
	m := &Manidator{
		Interval: DefaultInterval,
		done:     make(chan struct{}),
		cancel:   make(chan struct{}),
	}

	m.Add(dators...)

	return m
}

// Add includes one or more additional Dators in the Manidator's output.
func (this *Manidator) Add(dators ...Dator) {
	this.dators = append(this.dators, dators...)
}

// init determines the maximum width of the name field for all Dators and builds a format string based on that width.
func (this *Manidator) init() {
	for _, d := range this.dators {
		l := len(d.GetName())
		if l > this.nameWidth {
			this.nameWidth = l
		}
	}

	this.fmtString = fmt.Sprintf("%%%ds> %%s\n", this.nameWidth)
}

// Begin launches the Manidator and displays the output from each registered Dator. The loop runs in its own goroutine,
// so this function returns immediately.
func (this *Manidator) Begin(ctx context.Context) {
	this.wg.Add(1)
	this.init()

	go func() {
		defer this.wg.Done()

		for {
			this.printLines()

			select {
			case <-time.After(this.Interval):
				this.eraseLines()
			case <-this.Done():
				// all Dators finished
				return
			case <-this.cancel:
				// manidator was stopped
				return
			case <-ctx.Done():
				// app was stopped
				return
			}
		}
	}()
}

// Done allows the Manidator owner to know when all Dators have finished providing output.
func (this *Manidator) Done() chan struct{} {
	return this.done
}

// Stop allows the Manidator owner to request that all output come to a halt.
func (this *Manidator) Stop() {
	close(this.cancel)
	this.wg.Wait()
}

// printLines iterates over each Dator, printing the last line of output from each one.
func (this *Manidator) printLines() {
	closedCnt := 0

	width, _, err := terminal.GetSize(0)
	if err != nil {
		width = DefaultTerminalWidth
	}

	// Determine the max length for each line of output. This is computed on each pass to account for the user resizing
	// their terminal.
	maxWidth := width - this.nameWidth - 6

	for _, dator := range this.dators {
		name := dator.GetName()

		// just in case the name changed after Begin was called
		if len(name) > this.nameWidth {
			name = name[this.nameWidth:]
		}

		line := dator.GetLastLine()

		// ellide lines that are longer than the max width of the terminal
		if len(line) > maxWidth {
			line = "..." + line[len(line)-maxWidth:]
		}

		fmt.Printf(this.fmtString, name, line)

		if dator.IsClosed() {
			closedCnt++
		}
	}

	if closedCnt == len(this.dators) {
		close(this.done)
	}
}

// eraseLines wipes out each line of output from each Dator.
func (this *Manidator) eraseLines() {
	fmt.Printf(strings.Repeat(eraseLine, len(this.dators)))
}
