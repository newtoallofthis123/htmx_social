package main

import "fmt"

func main() {
	fmt.Println("Hello World")

	store, err := NewStore(GetEnv())

	// Panic ID: 1
	if err != nil {
		panic(fmt.Sprintf("Error creating store (Panic ID: 1): %s", err))
	}

	listenAddr := "localhost:2468"

	apiServer := NewApiServer(store, listenAddr).Start()

	// Panic ID: 2
	if err := apiServer.Run(listenAddr); err != nil {
		panic(fmt.Sprintf("Error with server (Panic ID: 2): %s", err))
	}
}
