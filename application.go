package main

func main() {
	a := App{}
	a.Initialize("sa", "Password95", "M", "")
	a.Run(":8080")
}
