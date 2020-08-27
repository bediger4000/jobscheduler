package main

/*
 * Create a heap.Heap from the command line, use that heap to
 * exercize heap.Draw, which puts GraphViz dot format output
 * on stdout.
 *
 * type IntegerNode holds integers, matches interface heap.Node.
 */

import (
	"fmt"
	"jobscheduler/heap"
	"os"
	"strconv"
)

type IntegerNode struct {
	Data int64
}

func (n *IntegerNode) Value() int64 {
	return n.Data
}

func (n *IntegerNode) IsNil() bool {
	return n == nil
}

func (n *IntegerNode) String() string {
	return fmt.Sprintf("%d", n.Data)
}

func main() {
	var h heap.Heap

	for _, str := range os.Args[1:] {
		n, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			h = h.Insert(&IntegerNode{Data: n})
		}
	}

	heap.Draw(h)
}
