/*
 * @Date: 2023-07-20 10:24:45
 * @LastEditTime: 2023-07-20 10:24:57
 * @Description:
 */
package xticker

import "time"

type Fn func() error
type MyTicker struct {
	MyTick *time.Ticker
	Runner Fn
	Done   chan bool //  done channel
}

// NewTicker ...
func NewTicker(duration int, f Fn) *MyTicker {
	return &MyTicker{
		MyTick: time.NewTicker(time.Duration(duration) * time.Second),
		Runner: f,
		Done:   make(chan bool),
	}
}

// Start ...
func (t *MyTicker) Start() {
	go func() {
		for {
			select {
			case <-t.MyTick.C:
				t.Runner()
			case <-t.Done:
				t.MyTick.Stop()
			}
		}
	}()
}

// Stop ...
func (t *MyTicker) Stop() {
	t.Done <- true
}
