package model

type (
	//StackHook represents hooks associated to a task
	StackHook struct {
		//Deploy specifies the hook tasks to run when a stack is deployed
		Deploy Hook
	}
)

//HasTasks returns true if the hook contains at least one task reference
func (r StackHook) HasTasks() bool {
	return r.Deploy.HasTasks()
}

func (r *StackHook) customize(with StackHook) error {
	return r.Deploy.customize(with.Deploy)
}

func (r StackHook) validate() ValidationErrors {
	return ErrorOnInvalid(r.Deploy)
}
