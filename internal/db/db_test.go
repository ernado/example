package db

func (suite *DBTestSuite) TestCreateTask() {
	ctx := suite.T().Context()
	task, err := suite.db.CreateTask(ctx, "task1")
	suite.NoError(err)

	tasks, err := suite.db.ListTasks(ctx)
	suite.NoError(err)
	suite.Len(tasks, 1)
	suite.Equal(task, tasks[0])
}
