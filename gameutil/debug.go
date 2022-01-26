package gameutil

import (
	"github.com/sirupsen/logrus"
	"os"
	"runtime/debug"
)

func PrintStack() {
	//打印调用堆
	stack := debug.Stack()
	logrus.WithFields(logrus.Fields{
		"stack": string(stack),
	}).Error("Print Stack")
	os.Stderr.Write(stack)
}
