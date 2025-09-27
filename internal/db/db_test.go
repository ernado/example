package entdb

import (
	"github.com/ernado/example"
)

func (suite *DBTestSuite) TestCreateTask() {
	ctx := suite.T().Context()
	err := suite.db.CreateTask(ctx, &example.Task{Title: "task1"})
	suite.NoError(err)

	tasks, err := suite.db.ListTasks(ctx)
	suite.NoError(err)
	suite.Len(tasks, 1)
	suite.Equal("task1", tasks[0].Title)
}
