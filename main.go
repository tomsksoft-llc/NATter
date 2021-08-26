package main

import (
	"NATter/cmd"
	"NATter/log"
)

func main() {
	natter := cmd.NewNATter()

	if err := natter.Run(); err != nil {
		log.Fatal(err)
	}
}
