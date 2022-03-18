package utils

import (
	"fmt"
	"runtime"
)

func EP(err error) error {
	if err != nil {
		fmt.Println(runtime.Caller(1))
		fmt.Println("error:", err.Error())
	}
	return err
}
