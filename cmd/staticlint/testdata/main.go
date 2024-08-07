package main

import "os"

func main() {
	os.Exit(0) // want "direct call to os.Exit found, consider returning an error instead"
}
