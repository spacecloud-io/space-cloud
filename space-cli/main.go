package main

import (
	"fmt"
	"os"
)

func main() {
	rootCmd, err := getmodule()
	err = rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
