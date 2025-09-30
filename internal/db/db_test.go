package db

func (suite *DBTestSuite) TestCreateTask() {
	ctx := suite.T().Context()

	task, err := suite.db.CreateTask(ctx, "Test Task")
	suite.Require().NoError(err)
	suite.Require().NotNil(task)

	tasks, err := suite.db.ListTasks(ctx)
	suite.Require().NoError(err)
	suite.Require().Len(tasks, 1)
	suite.Require().Equal(task, tasks[0])

	err = suite.db.DeleteTask(ctx, task.ID)
	suite.Require().NoError(err)

	tasks, err = suite.db.ListTasks(ctx)
	suite.Require().NoError(err)
	suite.Require().Len(tasks, 0)
}
