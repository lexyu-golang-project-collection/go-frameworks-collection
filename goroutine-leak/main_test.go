package main

import (
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestLeak(t *testing.T) {
	Leak()
}

func TestLeakWithGoleak(t *testing.T) {
	defer goleak.VerifyNone(t)
	Leak()
}
