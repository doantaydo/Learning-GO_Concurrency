package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	fmt.Println("Test main")
	stdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	_ = w.Close()
	result, _ := io.ReadAll(r)
	output := string(result)
	os.Stdout = stdOut

	if !strings.Contains(output, "$34320.00") {
		t.Error("Wrong balance returned")
	}
}
