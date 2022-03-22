package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/wqvoon/cbox/pkg/log"
)

// TODO: 用来展示一个勉强能看的进度条，后面再做改进，NewTask 的 fn 不应该有任何标准输出/标准错误

type Task struct {
	hint string
	fn   func()

	finished bool
	wg       *sync.WaitGroup
	done     chan struct{}
}

func NewTask(hint string, fn func()) *Task {
	return &Task{
		hint: hint,
		fn:   fn,
		wg:   new(sync.WaitGroup),
		done: make(chan struct{}),
	}
}

func (t *Task) Start() {
	if t.finished {
		log.Errorln("can not Start a finished task")
	}

	t.wg.Add(1)
	fmt.Print(t.hint)

	go func() {
		for {
			select {
			case <-t.done:
				fmt.Println("done")
				t.wg.Done()
				return
			default:
				fmt.Print(".")
				time.Sleep(time.Second)
			}
		}
	}()

	t.fn()
	t.Done()
}

func (t *Task) Done() {
	if t.finished {
		return
	}

	t.done <- struct{}{}
	t.wg.Wait()

	close(t.done)
	t.finished = true
}
