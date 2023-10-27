package main

import (
	"fmt"
)

func ErrorHandler(e error) {
	if e != nil {
		fmt.Println("Unsolved error!",e)
	}
}

/*
处理开头多余回车与字符串
*/
func clsStr(s *string){ 
	if *s != "" {
		for {
			if (*s)[0] == '\n' || (*s)[0] == ' '{
				*s = (*s)[1:]
			} else {
				break
			}
		}
	}
}