package statsd

import (
	"bytes"
	"testing"
	"time"
)

func assert(t *testing.T, value, control string) {
	if value != control {
		t.Errorf("incorrect command, want '%s', got '%s'", control, value)
	}
}

func TestPrefix(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewClient(buf)
	c.Prefix("foo.bar.baz.")
	err := c.Increment("incr", 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	c.Flush()
	assert(t, buf.String(), "foo.bar.baz.incr:1|c")
}

func TestIncrement(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewClient(buf)
	err := c.Increment("incr", 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	c.Flush()
	assert(t, buf.String(), "incr:1|c")
}

func TestIncr(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewClient(buf)
	err := c.Incr("incr")
	if err != nil {
		t.Fatal(err)
	}
	c.Flush()
	assert(t, buf.String(), "incr:1|c")
}

func TestDecrement(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewClient(buf)
	err := c.Decrement("decr", 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	c.Flush()
	assert(t, buf.String(), "decr:-1|c")
}

func TestDecr(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewClient(buf)
	err := c.Decr("decr")
	if err != nil {
		t.Fatal(err)
	}
	c.Flush()
	assert(t, buf.String(), "decr:-1|c")
}

func TestDuration(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewClient(buf)
	err := c.Duration("timing", time.Duration(123456789), 1)
	if err != nil {
		t.Fatal(err)
	}
	c.Flush()
	assert(t, buf.String(), "timing:123|ms")
}

func TestIncrementRate(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewClient(buf)
	err := c.Increment("incr", 1, 0.99)
	if err != nil {
		t.Fatal(err)
	}
	c.Flush()
	assert(t, buf.String(), "incr:1|c|@0.99")
}

func TestPreciseRate(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewClient(buf)
	// The real use case here is rates like 0.0001.
	err := c.Increment("incr", 1, 0.99901)
	if err != nil {
		t.Fatal(err)
	}
	c.Flush()
	assert(t, buf.String(), "incr:1|c|@0.99901")
}

func TestRate(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewClient(buf)
	err := c.Increment("incr", 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	c.Flush()
	assert(t, buf.String(), "")
}

func TestGauge(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewClient(buf)
	err := c.Gauge("gauge", 300, 1)
	if err != nil {
		t.Fatal(err)
	}
	c.Flush()
	assert(t, buf.String(), "gauge:300|g")
}

func TestIncrementGauge(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewClient(buf)
	err := c.IncrementGauge("gauge", 10, 1)
	if err != nil {
		t.Fatal(err)
	}
	c.Flush()
	assert(t, buf.String(), "gauge:+10|g")
}

func TestDecrementGauge(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewClient(buf)
	err := c.DecrementGauge("gauge", 4, 1)
	if err != nil {
		t.Fatal(err)
	}
	c.Flush()
	assert(t, buf.String(), "gauge:-4|g")
}

func TestUnique(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewClient(buf)
	err := c.Unique("unique", 765, 1)
	if err != nil {
		t.Fatal(err)
	}
	c.Flush()
	assert(t, buf.String(), "unique:765|s")
}

var millisecondTests = []struct {
	duration time.Duration
	control  int
}{
	{
		duration: 350 * time.Millisecond,
		control:  350,
	},
	{
		duration: 5 * time.Second,
		control:  5000,
	},
	{
		duration: 50 * time.Nanosecond,
		control:  0,
	},
}

func TestMilliseconds(t *testing.T) {
	for i, mt := range millisecondTests {
		value := millisecond(mt.duration)
		if value != mt.control {
			t.Errorf("%d: incorrect value, want %d, got %d", i, mt.control, value)
		}
	}
}

func TestTiming(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewClient(buf)
	err := c.Timing("timing", 350, 1)
	if err != nil {
		t.Fatal(err)
	}
	c.Flush()
	assert(t, buf.String(), "timing:350|ms")
}

func TestTime(t *testing.T) {
	buf := new(bytes.Buffer)
	c := NewClient(buf)
	err := c.Time("time", 1, func() { time.Sleep(50e6) })
	if err != nil {
		t.Fatal(err)
	}
}
