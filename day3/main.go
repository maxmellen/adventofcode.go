package main

import (
	"fmt"
	"io"
	"os"
)

type claim struct {
	id   uint
	x, y uint
	w, h uint
}

func main() {
	file := os.Stdin
	parsedClaims := make(chan claim)
	go parseClaims(file, parsedClaims)
	intactClaim := filterClaims(parsedClaims)
	fmt.Println(intactClaim.id)
}

func parseClaims(r io.Reader, out chan<- claim) {
	for {
		c := claim{}

		_, err := fmt.Fscanf(r, "#%d @ %d,%d: %dx%d\n",
			&c.id, &c.x, &c.y, &c.w, &c.h)

		if err != nil && err != io.ErrUnexpectedEOF {
			fmt.Println("Could not parse claims:", err)
			os.Exit(1)
		}

		if err == io.ErrUnexpectedEOF {
			close(out)
			break
		}

		out <- c
	}
}

func filterClaims(in <-chan claim) claim {
	entryPoint := make(chan claim, 1)
	ch1 := entryPoint

	for c := range in {
		entryPoint <- c
		ch2 := make(chan claim)
		go doFilterClaims(ch1, ch2, c)
		ch1 = ch2
	}

	return <-ch1
}

func doFilterClaims(in <-chan claim, out chan<- claim, filter claim) {
	for c := range in {
		if c.id == filter.id || !overlap(c, filter) {
			out <- c
		}
	}
}

func overlap(c1 claim, c2 claim) bool {
	c1left := c1.x
	c1right := c1.x + c1.w
	c1top := c1.y
	c1bottom := c1.y + c1.h

	c2left := c2.x
	c2right := c2.x + c2.w
	c2top := c2.y
	c2bottom := c2.y + c2.h

	return !(c1left < c2left && c1right <= c2left ||
		c1left >= c2right && c1right > c2right ||
		c1top < c2top && c1bottom <= c2top ||
		c1top >= c2bottom && c1bottom > c2bottom)
}
