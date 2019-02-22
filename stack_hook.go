package model

type (
	//StackHook represents hooks associated to a task
	StackHook struct {
		//Deploy specifies the hook tasks to run when a stack is deployed
		Deploy Hook
		//Undeploy specifies the hook tasks to run when a stack is undeployed
		Undeploy Hook
	}
)

//HasTasks returns true if the hook contains at least one task reference
func (r StackHook) HasTasks() bool {
	return r.Deploy.HasTasks() ||
		r.Undeploy.HasTasks()
}

func (r *StackHook) merge(other StackHook) error {
	if err := r.Deploy.merge(other.Deploy); err != nil {
		return err
	}
	if err := r.Undeploy.merge(other.Undeploy); err != nil {
		return err
	}

	return nil
}

func (r StackHook) validate() ValidationErrors {
	return ErrorOnInvalid(r.Deploy, r.Undeploy)
}
