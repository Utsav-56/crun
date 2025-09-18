package main

import "fmt"

type Ulog struct {
	count int
}

func (l *Ulog) println(format string, args ...interface{}) {
	l.count++
	if len(args) == 0 {
		fmt.Println(format)
	} else {
		fmt.Printf(format+"\n", args...)
	}
}

func (l *Ulog) clear() {
	if l.count > 0 && !flags.verbose {
		clearLastLines(l.count)
		l.count = 0
	}
}
