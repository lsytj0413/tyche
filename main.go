package main

import (
	"fmt"

	"github.com/lsytj0413/tyche/tcb"
)

func main() {
	awards, err := tcb.Fetch()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", awards)
	fmt.Println("tyche")
}
