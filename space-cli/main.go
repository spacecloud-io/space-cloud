package main

import (
	"fmt"
	"os"
)

func main() {
	rootCmd, err := getmodule()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
