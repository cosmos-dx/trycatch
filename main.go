package trycatch

import (
	"fmt"
	"log"
	"os"
	"time"
)

type TryCatch struct {
	work         func() error
	err          error
	defaultCatch func(error)
	finalize     func()
	retries      int
	retryDelay   time.Duration
	logger       *log.Logger
	metrics      func(error)
}

func Try(work func() error) *TryCatch {
	return &TryCatch{
		work:       work,
		retries:    0,
		retryDelay: 0,
		logger:     log.New(os.Stdout, "[TRYCATCH] ", log.LstdFlags),
	}
}

func (t *TryCatch) Catch(handler func(error)) *TryCatch {
	t.defaultCatch = handler
	return t
}

func (t *TryCatch) Finally(finalize func()) *TryCatch {
	t.finalize = finalize
	return t
}

func (t *TryCatch) Recover() *TryCatch {
	defer func() {
		if r := recover(); r != nil {
			t.err = fmt.Errorf("panic recovered: %v", r)
			t.logger.Printf("Panic recovered: %v", r)
			t.Execute()
		}
	}()
	t.executeWork() // Call work function
	return t
}

func (t *TryCatch) executeWork() {
	t.err = t.work()
	if t.err != nil {
		t.Execute()
	}
}

func (t *TryCatch) Execute() {
	if t.err != nil {
		if t.defaultCatch != nil {
			t.defaultCatch(t.err)
		} else {
			t.logger.Printf("Unhandled error: %v", t.err)
		}
		if t.metrics != nil {
			t.metrics(t.err)
		}
	}
	if t.finalize != nil {
		t.finalize()
	}
}
