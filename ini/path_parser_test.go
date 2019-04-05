package inilib

import (
	"fmt"
	"testing"
)

var path string = "[].server"

func Test_parsePath(t *testing.T) {
	paths := parsePath(path)
	if len(paths) <= 0 {
		t.Error("pase path error")
	}

	for _, value := range paths {
		fmt.Println(value)
	}

	t.Log("test ok!")
}
