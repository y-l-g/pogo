package supervisor

import (
	"sync"
)

type WaitGroup struct {
	Wg          sync.WaitGroup
	OwnerPoolID int64
}

func (wg *WaitGroup) Add(delta int64) { wg.Wg.Add(int(delta)) }
func (wg *WaitGroup) Done()           { wg.Wg.Done() }
func (wg *WaitGroup) Wait()           { wg.Wg.Wait() }

type Channel struct {
	Ch          chan string
	OwnerPoolID int64
}

func (c *Channel) Init(capacity int64) { c.Ch = make(chan string, int(capacity)) }
func (c *Channel) Push(value string)   { c.Ch <- value }
func (c *Channel) Pop() string {
	val, ok := <-c.Ch
	if !ok {
		return ""
	}
	return val
}
func (c *Channel) Close() { close(c.Ch) }
