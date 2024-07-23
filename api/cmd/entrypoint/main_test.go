package main

import (
	"api/cmd/config"
	"testing"
)

func Test_startMain(t *testing.T) {
	config.Init()
	startMain()
}
