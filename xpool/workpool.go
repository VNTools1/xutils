// Package xpool
/*
 * @Date: 2023-07-11 13:55:42
 * @LastEditTime: 2023-07-12 14:28:40
 * @Description:
 */
package xpool

import (
	"context"
	"fmt"
	"sync"
)

type TaskFunc func(args ...interface{})

type Task struct {
	F    TaskFunc
	Args interface{}
}

type WorkPool struct {
	Pool       chan *Task      //定义任务池
	WorkCount  int             //工作线程数量,决定初始化几个goroutine
	StopCtx    context.Context //上下文
	StopCancel context.CancelFunc
	WG         sync.WaitGroup //阻塞计数器
}

// Execute ...
func (t *Task) Execute(args ...interface{}) {
	t.F(args...)
}

// NewWorkPool ...
func NewWorkPool(workerCount int, len int) *WorkPool {
	return &WorkPool{
		WorkCount: workerCount,
		Pool:      make(chan *Task, len),
	}
}

// PushTask ...
func (w *WorkPool) PushTask(task *Task) {
	w.Pool <- task
}

// Work ...
func (w *WorkPool) Work(wid int) {
	for {
		select {
		case <-w.StopCtx.Done():
			w.WG.Done()
			fmt.Printf("线程%d 退出执行了 \n", wid)
			return
		case t := <-w.Pool:
			if t != nil {
				t.Execute()
				fmt.Printf("f被线程%d执行了，参数为%v \n", wid, t.Args)
			}

		}
	}

}

// Start ...
func (w *WorkPool) Start() *WorkPool {
	//定义好worker数量
	w.WG.Add(w.WorkCount)
	w.StopCtx, w.StopCancel = context.WithCancel(context.Background())
	for i := 0; i < w.WorkCount; i++ {
		//定义多少个协程来工作
		go w.Work(i)
	}
	return w
}

// Stop ...
func (w *WorkPool) Stop() {
	w.StopCancel()
	w.WG.Wait()
}
