package main

func main() {
	a := App{}
	a.Initialize("sa", "Password95", "M.cykigjdaqb15.us-east-1.rds.amazonaws.com", "3306", "M", "")
	a.Run(":2345")
}
