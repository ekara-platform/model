package model

type (
	//NodeHook represents hooks associated to a node set
	NodeHook struct {
		//Provisioned specifies the hook tasks to run when a node set is provisioned
		Provision Hook
	}
)

//HasTasks returns true if the hook contains at least one task reference
func (r NodeHook) HasTasks() bool {
	return r.Provision.HasTasks()
}

func (r *NodeHook) customize(with NodeHook) error {
	return r.Provision.customize(with.Provision)
}

func (r NodeHook) validate() ValidationErrors {
	return ErrorOnInvalid(r.Provision)
}
