package main

import "github.com/cleodora-forecasting/cleodora/cleosrv/cleosrv"

func main() {
	if err := cleosrv.Start(); err != nil {
		panic(err)
	}
}
