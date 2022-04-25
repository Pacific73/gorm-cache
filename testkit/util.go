package testkit

import (
	"fmt"
	"time"
)

func log(format string, a ...interface{}) {
	timeStr := time.Now().Format("2006-01-02 15:04:05.999")
	fmt.Printf(timeStr+" "+format+"\n", a...)
}

func timer(name string, f func() error) error {
	start := time.Now()
	fmt.Printf("[%s] start ...\n", name)
	err := f()
	duration := time.Now().Sub(start)
	fmt.Printf("[%s] finished. cost: %.3fs\n", name, duration.Seconds())
	return err
}
