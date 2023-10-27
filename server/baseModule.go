package main

import (
	"fmt"
)

func ErrorHandler(e error) {
	if e != nil {
		fmt.Println("Unsolved error!",e)
	}
}