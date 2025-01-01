package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_updateMessage(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	updateMessage("Test updateMessage", &wg)
	wg.Wait()

	if !strings.Contains(msg, "Test updateMessage") {
		t.Errorf("MSG: Expected to find 'Test updateMessage' but it is not there")
	}
}

func Test_printMessage(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	msg = "Test printMessage"
	printMessage()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "Test printMessage") {
		t.Errorf("STDOUT: Expected to find 'Test printMessage' but it is not there")
	}
}

func Test_main(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "Hello, universe!") {
		t.Errorf("STDOUT: Expected to find 'Hello, universe!' but it is not there")
	}
	if !strings.Contains(output, "Hello, cosmos!") {
		t.Errorf("STDOUT: Expected to find 'Hello, cosmos!' but it is not there")
	}
	if !strings.Contains(output, "Hello, world!") {
		t.Errorf("STDOUT: Expected to find 'Hello, world!' but it is not there")
	}
	if !strings.Contains(msg, "Hello, world!") {
		t.Errorf("MSG: Expected to find 'Hello, world!' but it is not there")
	}
}
