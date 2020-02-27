package model

type (
	//NodeHook represents hooks associated to a node set
	NodeHook struct {
		//Create specifies the hook tasks to run when a node set is created
		Create Hook
		//Destroy specifies the hook tasks to run when a node set is destroyed
		Destroy Hook
	}
)

//HasTasks returns true if the hook contains at least one task reference
func (r NodeHook) HasTasks() bool {
	return r.Create.HasTasks() || r.Destroy.HasTasks()
}

func (r *NodeHook) customize(with NodeHook) error {
	err := r.Create.customize(with.Create)
	if err != nil {
		return err
	}
	err = r.Destroy.customize(with.Destroy)
	if err != nil {
		return err
	}
	return nil
}

func (r NodeHook) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	vErrs.merge(ErrorOnInvalid(r.Create))
	vErrs.merge(ErrorOnInvalid(r.Destroy))
	return vErrs
}
