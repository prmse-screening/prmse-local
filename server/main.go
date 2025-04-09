package main

import (
	"server/internal/config"
	"server/internal/data"
	"server/internal/logger"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	config.Init()
	logger.Init()
	_, _ = data.NewDatabase()

}
