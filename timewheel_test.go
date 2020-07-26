package timewheel

import (
	"testing"
	"time"
)

func TestAddTimer(t *testing.T) {
	tw := NewTimeWheel(time.Second)
	tw.AddTimer(2*time.Second, func(arg interface{}){
		beforeTime := arg.(time.Time)
		nowTime := time.Now()
		actual := int(nowTime.Sub(beforeTime) / time.Second)
		expected := 2
		if actual != expected {
			t.Errorf("add timer test failed. actual:%d expected:%d", actual, expected)
		}
	}, time.Now())
	time.Sleep(3*time.Second)
}

func TestAddTimerDeadline(t *testing.T) {
	tw := NewTimeWheel(time.Second)
	tw.AddTimerDeadline(time.Now().Add(2*time.Second), func(arg interface{}){
		beforeTime := arg.(time.Time)
		nowTime := time.Now()
		actual := int(nowTime.Sub(beforeTime) / time.Second)
		expected := 2
		if actual != expected {
			t.Errorf("add timer test failed. actual:%d expected:%d", actual, expected)
		}
	}, time.Now())
	time.Sleep(3*time.Second)
}

func TestAddTimerTick(t *testing.T) {
	tw := NewTimeWheel(time.Second)
	tw.AddTimerTick(2, func(arg interface{}){
		beforeTime := arg.(time.Time)
		nowTime := time.Now()
		actual := int(nowTime.Sub(beforeTime) / time.Second)
		expected := 2
		if actual != expected {
			t.Errorf("add timer test failed. actual:%d expected:%d", actual, expected)
		}
	}, time.Now())
	time.Sleep(3*time.Second)
}

func TestAddTimerTickDeadline(t *testing.T) {
	tw := NewTimeWheel(time.Second)
	tw.AddTimerTickDeadline(tw.GetTick(time.Now())+2, func(arg interface{}){
		beforeTime := arg.(time.Time)
		nowTime := time.Now()
		actual := int(nowTime.Sub(beforeTime) / time.Second)
		expected := 2
		if actual != expected {
			t.Errorf("add timer test failed. actual:%d expected:%d", actual, expected)
		}
	}, time.Now())
	time.Sleep(3*time.Second)
}



