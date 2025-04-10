package main

import (
	"fmt"
	"server/internal/config"
	"server/internal/data"
	"server/internal/data/db"
	"server/internal/logger"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	config.Init()
	logger.Init()
	database, _ := data.NewDatabase()
	repo := db.NewTasksRepo(database)
	task, err := repo.NextTask()
	fmt.Println(task)
	fmt.Println(err != nil)
}
