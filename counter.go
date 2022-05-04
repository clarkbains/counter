package counter

import (
	"context"
	"errors"
	"fmt"
	"time"
)
type Counter struct{
	value int
	timeout time.Duration
}

//Create a new counter
func NewCounter () Counter {
	c := Counter{value:0, timeout: 500 * time.Millisecond}
	return c
}

//Log the current value of the counter
func (c Counter) LogValue () {
	fmt.Printf("The new value is %d\n", c.value)
}

//Add one to the counter
func (c *Counter) AddOne () {
	c.value += 1
}

//Set delay for adding with context, before the add is completed
func (c *Counter) SetAddDelay (d time.Duration) {
	c.timeout = d
}

//Add one, but wait for the time set by SetAddDelay, returning error if context expires first.
func (c *Counter) AddOneWithContext(ctx context.Context) error {
	type addRes struct {
		finished bool
		newValue int
	}
	channel := make(chan addRes)
	
	go func (ctx context.Context)  {
		select {
		case <- ctx.Done():
			channel <- addRes{}
		case <- time.After(c.timeout):
			c.AddOne()
			channel <- addRes{finished: true, newValue: c.value}
		}
		close(channel)
	}(ctx)

	res := <- channel
	if !res.finished {
		return errors.New("context cancelled")
	}
	return nil
}