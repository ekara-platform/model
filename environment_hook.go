package model

type (
	//EnvironmentHooks represents hooks associated to the environment
	EnvironmentHooks struct {
		//Create specifies the hook tasks to run when the environment is created
		Create Hook
		//Deploy specifies the hook tasks to run at the environment deployment
		Deploy Hook
	}
)

//HasTasks returns true if the hook contains at least one task reference
func (r EnvironmentHooks) HasTasks() bool {
	return r.Create.HasTasks() ||
		r.Deploy.HasTasks()
}

func (r *EnvironmentHooks) customize(with EnvironmentHooks) error {
	if err := r.Create.customize(with.Create); err != nil {
		return err
	}
	return r.Deploy.customize(with.Deploy)
}

func (r EnvironmentHooks) validate() ValidationErrors {
	return ErrorOnInvalid(r.Create, r.Deploy)
}
