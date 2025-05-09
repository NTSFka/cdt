package main

import (
	. "cdt/pkg"
	"fmt"
)

func main() {
	if err := RunMain(NewRunContext()); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}
