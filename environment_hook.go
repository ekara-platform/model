package model

type (
	//EnvironmentHooks represents hooks associated to the environment
	EnvironmentHooks struct {
		//Provisione specifies the hook tasks to run when the environment is provisioned
		Provision Hook
		//Deploy specifies the hook tasks to run at the environment deployment
		Deploy Hook
	}
)

//HasTasks returns true if the hook contains at least one task reference
func (r EnvironmentHooks) HasTasks() bool {
	return r.Provision.HasTasks() ||
		r.Deploy.HasTasks()
}

func (r *EnvironmentHooks) merge(other EnvironmentHooks) error {
	if err := r.Provision.merge(other.Provision); err != nil {
		return err
	}
	return r.Deploy.merge(other.Deploy)
}

func (r EnvironmentHooks) validate() ValidationErrors {
	return ErrorOnInvalid(r.Provision, r.Deploy)
}
