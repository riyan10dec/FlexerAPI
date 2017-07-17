package main

import model "FlexerAPI/Model"

func main() {
	var Config model.Config
	Config.LoadConfiguration("config.json")

	a := App{Config: Config}
	a.Initialize()
	a.Run(":2345")
}
