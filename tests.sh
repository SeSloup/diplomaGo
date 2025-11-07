#/bin/bash

go test -count=1 -run ^TestNextDate$ ./tests
go test -count=1 -run ^TestAddTask$ ./tests
go test -count=1 -run ^TestTasks$ ./tests
go test -count=1 -run ^TestEditTask$ ./tests
go test -count=1 -run ^TestDone$ ./tests
go test -count=1 -run ^TestDelTask$ ./tests
