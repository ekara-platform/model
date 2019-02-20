package model

type (
	//EnvironmentHooks represents hooks associated to the environment
	EnvironmentHooks struct {
		//Provisione specifies the hook tasks to run when the environment is provisioned
		Provision Hook
		//Deploy specifies the hook tasks to run at the environment deployment
		Deploy Hook
		//Undeploy specifies the hook tasks to run when the environment is undeployed
		Undeploy Hook
		//Destroy specifies the hook tasks to run when the environment is destroyed
		Destroy Hook
	}
)

//HasTasks returns true if the hook contains at least one task reference
func (r EnvironmentHooks) HasTasks() bool {
	return r.Provision.HasTasks() ||
		r.Deploy.HasTasks() ||
		r.Undeploy.HasTasks() ||
		r.Destroy.HasTasks()
}

func (r *EnvironmentHooks) merge(other EnvironmentHooks) error {
	if err := r.Provision.merge(other.Provision); err != nil {
		return err
	}
	if err := r.Deploy.merge(other.Deploy); err != nil {
		return err
	}
	if err := r.Undeploy.merge(other.Undeploy); err != nil {
		return err
	}
	if err := r.Destroy.merge(other.Destroy); err != nil {
		return err
	}
	return nil
}

func (r EnvironmentHooks) validate() ValidationErrors {
	return ErrorOnInvalid(r.Provision, r.Deploy, r.Undeploy, r.Destroy)
}
