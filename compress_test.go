package glog

import (
	"flag"
	"testing"
)

func TestRunMain(t *testing.T) {
	flag.Parse()

	Info("nothing")

	Flush()
}
