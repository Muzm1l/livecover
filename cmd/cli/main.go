package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("coverage CLI scaffold")
	fmt.Println("usage:")
	fmt.Println("  agent      run host agent")
	fmt.Println("  aggregator run central aggregator")
	fmt.Println("  web        run web ui placeholder")
	os.Exit(0)
}
