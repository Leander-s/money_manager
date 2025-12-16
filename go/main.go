package main

func main() {
	ctx := initContext()
	runServer(ctx)
	deinitContext(ctx)
}
