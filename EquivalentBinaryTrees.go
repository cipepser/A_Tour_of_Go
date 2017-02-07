package main

import (
	"golang.org/x/tour/tree"
	"fmt"
)

func Walk(t *tree.Tree, ch chan int) {
	if t.Left != nil {
		Walk(t.Left, ch)
	}

	ch <- t.Value

	if t.Right != nil {
		Walk(t.Right, ch)
	}		
}

func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	
	go func () {
		Walk(t1, ch1)
		close(ch1)
	} ()
	
	go func () {
		Walk(t2, ch2)
		close(ch2)
	} ()
	
	for {
		v1, ok1 := <- ch1
		v2, ok2 := <- ch2
		
		if !ok1 || !ok2 {
			break
		}

		if v1 != v2 {
			return false
		}
	}
	
	return true
	
	
}

func main() {
	t := tree.New(1)
	t = tree.New(1)
	ch := make(chan int)
	
	go func () {
		Walk(t, ch)
		close(ch)
	} ()
	
	for i := range ch {
		fmt.Println(i)
	}
	
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}