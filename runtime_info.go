package model

type (
	//RunTimeInfo the context passed to all ekara templates
	RunTimeInfo struct {
		//TargetType represents the type of the target supporting the current action
		TargetType string
		//TargetName represents the name of the target supporting the current action
		TargetName string
	}
)

func createRunTimeInfo() *RunTimeInfo {
	return &RunTimeInfo{}
}

//SetTarget defines the target of the running action
func (cc *RunTimeInfo) SetTarget(t Describable) {
	cc.TargetType = t.DescType()
	cc.TargetName = t.DescName()
}
