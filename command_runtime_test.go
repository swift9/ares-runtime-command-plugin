package runtime_test

import (
	"log"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	time.Sleep(1 * time.Hour)
}
