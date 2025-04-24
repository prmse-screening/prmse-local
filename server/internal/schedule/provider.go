package schedule

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewTasksScheduler, NewTaskCleaner)
