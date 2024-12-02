package main

import isuports "github.com/isucon/isucon12-qualify/webapp/go"

func main() {
	if err := isuports.CleanDB(); err != nil {
		panic(err)
	}
}
