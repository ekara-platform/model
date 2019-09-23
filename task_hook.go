package model

type (
	//TaskHook represents hooks associated to a task
	TaskHook struct {
		//Execute specifies the hook tasks to run when a task is executed
		Execute Hook
	}
)

//HasTasks returns true if the hook contains at least one task reference
func (r TaskHook) HasTasks() bool {
	return r.Execute.HasTasks()
}

func (r *TaskHook) customize(with TaskHook) error {
	return r.Execute.customize(with.Execute)
}

func (r TaskHook) validate() ValidationErrors {
	return ErrorOnInvalid(r.Execute)
}
