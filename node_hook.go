package model

type (
	//NodeHook represents hooks associated to a node set
	NodeHook struct {
		//Create specifies the hook tasks to run when a node set is creaed
		Create Hook
	}
)

//HasTasks returns true if the hook contains at least one task reference
func (r NodeHook) HasTasks() bool {
	return r.Create.HasTasks()
}

func (r *NodeHook) customize(with NodeHook) error {
	return r.Create.customize(with.Create)
}

func (r NodeHook) validate() ValidationErrors {
	return ErrorOnInvalid(r.Create)
}
