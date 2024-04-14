package main

import (
	"fmt"
	"github.com/lookeme/short-url/internal/server/http"
)

func main() {
	err := http.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
