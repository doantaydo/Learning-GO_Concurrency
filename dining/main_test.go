package main

import (
	"testing"
	"time"
)

func Test_dine(t *testing.T) {
	eatTime = 0 * time.Second
	thinkTime = 0 * time.Second
	sleepTime = 0 * time.Second

	for i := 0; i < 10; i++ {
		orderFinishedName = []string{}
		dine()
		if len(orderFinishedName) != len(philosophers) {
			t.Errorf("At loop %d, incorrect length of slice; expected %d but got %d", i, len(philosophers), len(orderFinishedName))
		}
	}
}

func Test_dineWithVaryingDelay(t *testing.T) {
	var theTest = []struct {
		name  string
		delay time.Duration
	}{
		{"zero delay", time.Second * 0},
		{"quarter second delay", time.Millisecond * 25},
		{"half second delay", time.Millisecond * 50},
	}

	for _, e := range theTest {
		orderFinishedName = []string{}

		eatTime = e.delay
		thinkTime = e.delay
		sleepTime = e.delay

		dine()
		if len(orderFinishedName) != len(philosophers) {
			t.Errorf("%s: incorrect length of slice; expected %d but got %d", e.name, len(philosophers), len(orderFinishedName))
		}
	}
}
