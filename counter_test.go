package counter

import (
	"context"
	"testing"
	"time"
)

func TestInitial(t *testing.T) {
	c := NewCounter()
	if (c.value != 0) {
		t.Errorf("The initial value is not 0")
	}
	if (c.timeout != 500 * time.Millisecond) {
		t.Error("The initial delay is not 0")
	}
}

func TestBasic(t *testing.T) {
	c := NewCounter()
	for x := 0; x < 5; x++ {
		if x != c.value {
			t.Errorf("Value should be %d after %d increments, but it is %d instead.\n", x, x, c.value)
		}
		c.AddOne()
	}
}

func TestContextAdd(t *testing.T) {
	var expectedVal int = 0
	c := NewCounter()
	runTest := func (shouldTimeout bool, contextDeadline int, counterDelay int)  {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(contextDeadline * int(time.Millisecond)))
		defer cancel()
		c.SetAddDelay(time.Duration(counterDelay * int(time.Millisecond)))
		err := c.AddOneWithContext(ctx)

		if (c.value != expectedVal) {
			t.Errorf("Expected to get %d, got %d", expectedVal, c.value)
		}
		
		if err == nil && shouldTimeout {
			t.Errorf("Should have timed out, but didn't, with ctx deadline of %d, and counter delay of %d", contextDeadline, counterDelay)
		}

		if (err != nil){
			expectedVal++
		}
	}

	runTest(true, 200, 300)
}