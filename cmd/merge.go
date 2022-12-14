package main

import (
	"fmt"
)

func split(pin, pout chan int) {
	vall := <-pin
	valr, more := <-pin
	if more {
		// Spawn off recursive sorters
		lin := make(chan int)
		lout := make(chan int)
		go split(lout, lin)
		rin := make(chan int)
		rout := make(chan int)
		go split(rout, rin)

		// Send data from parent to recursive sorters
		lout <- vall
		rout <- valr
		for val := range pin {
			lout <- val
			lout, rout = rout, lout
		}
		close(lout)
		close(rout)

		// Merge sorted data received from children and send to parent
		merge(pout, lin, rin)

	} else {
		// Only one element to sort - just return it to parent
		pout <- vall
		close(pout)
	}
}

func merge(pout, lin, rin chan int) {
	// Receive and merge values from recursive sorters
	lval, lmore := <-lin
	rval, rmore := <-rin
	for rmore {
		if rval <= lval {
			pout <- rval
		} else {
			pout <- lval
			lval = rval
			lmore, rmore = rmore, lmore
			lin, rin = rin, lin
		}
		rval, rmore = <-rin
	}
	// Done with rin

	// Copy remaining values from lin to pout
	for lmore {
		pout <- lval
		lval, lmore = <-lin
	}
	close(pout)
}

func main() {
	input := [10]int{1, 4, 4, 2, 7, 3, 5, 8, 42, 0}

	in := make(chan int)
	out := make(chan int)

	go split(out, in)

	for _, val := range input {
		out <- val
	}
	close(out)

	for val := range in {
		fmt.Println(val)
	}
}
