package main

func main() {
	app := initServer()
	app.runServer()
	app.deInitServer()
}
