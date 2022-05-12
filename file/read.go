package main

import (
	"io/ioutil"
	"strings"
)

func Read(path string) (string, error) {
	if content, err := ioutil.ReadFile(path); err != nil {
		return "", err
	} else {
		return strings.TrimRight(string(content), "\r\n"), nil
	}
}
