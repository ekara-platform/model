package model

type (
	//EnvironmentHooks represents hooks associated to the environment
	EnvironmentHooks struct {
		//Init specifies the hook tasks to run before the environment is created
		Init Hook
		//Create specifies the hook tasks to run when the environment is created
		Create Hook
		//Install specifies the hook tasks to run at the environment installation
		Install Hook
		//Deploy specifies the hook tasks to run at the environment deployment
		Deploy Hook
		//Delete specifies the hook tasks to run at the environment deletion
		Delete Hook
	}
)

//HasTasks returns true if the hook contains at least one task reference
func (r EnvironmentHooks) HasTasks() bool {
	return r.Init.HasTasks() ||
		r.Create.HasTasks() ||
		r.Install.HasTasks() ||
		r.Deploy.HasTasks() ||
		r.Delete.HasTasks()
}

func (r *EnvironmentHooks) customize(with EnvironmentHooks) error {
	if err := r.Init.customize(with.Init); err != nil {
		return err
	}
	if err := r.Create.customize(with.Create); err != nil {
		return err
	}
	if err := r.Install.customize(with.Install); err != nil {
		return err
	}
	if err := r.Deploy.customize(with.Deploy); err != nil {
		return err
	}
	return r.Delete.customize(with.Delete)
}

func (r EnvironmentHooks) validate() ValidationErrors {
	return ErrorOnInvalid(r.Init, r.Create, r.Install, r.Deploy, r.Delete)
}
