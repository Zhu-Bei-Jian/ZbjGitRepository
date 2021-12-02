package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M){
   exitCode:=m.Run()

	// 退出
	os.Exit(exitCode)
}