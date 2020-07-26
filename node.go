package timewheel

import "sync"

var nodePool *sync.Pool

func init() {
	nodePool = &sync.Pool{
		New: func() interface{}{
			return &node{p:nodePool}
		},
	}
}

type node struct {
	expire int64  // 到期时间
	callback func(arg interface{})
	arg interface{}
	p *sync.Pool
}

func newNode(expire int64, callback func(arg interface{}), arg interface{}) *node {
	n := nodePool.Get().(*node)
	n.expire = expire
	n.callback = callback
	n.arg = arg
	return n
}

func (n *node) release() {
	if n != nil {
		n.expire = 0
		n.callback = nil
		n.arg = nil
		n.p.Put(n)
	}
}
