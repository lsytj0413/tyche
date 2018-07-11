package main

import (
	"fmt"

	"github.com/lsytj0413/tyche/tcb"
)

func main() {
	termList, err := tcb.FetchTermList()
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = termList

	_, err = tcb.FetchFromTerm(18077)
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Printf("%+v\n", awards)
	fmt.Println("tyche")
}
