package main

/*
 * Test package heap by creating a type that holds integers,
 * fits interface heap.Node, and sorts numbers given on the
 * command line.
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
	for len(h) > 0 {
		var n heap.Node
		h, n = h.Delete()
		fmt.Printf("%d\n", n.Value())
	}
}
