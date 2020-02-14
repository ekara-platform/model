package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var taskOtherT1, taskOtherT2, taskOriginT1, taskOriginT2 TaskRef

func TestTaskDescType(t *testing.T) {
	s := Task{}
	assert.Equal(t, s.DescType(), "Task")
}

func TestTaskDescName(t *testing.T) {
	s := Task{Name: "my_name"}
	assert.Equal(t, s.DescName(), "my_name")
}

func getTaskOrigin() *Task {
	t1 := &Task{
		Name:     "my_name",
		Playbook: "Playbook",
	}
	t1.EnvVars = make(map[string]string)
	t1.EnvVars["key1"] = "val1_target"
	t1.EnvVars["key2"] = "val2_target"
	t1.Parameters = make(map[string]interface{})
	t1.Parameters["key1"] = "val1_target"
	t1.Parameters["key2"] = "val2_target"

	t1.cRef = componentRef{
		ref: "cOriginal",
	}

	taskOriginT1 = TaskRef{ref: "T1"}
	taskOriginT2 = TaskRef{ref: "T2"}
	hT1 := TaskHook{}
	hT1.Execute.Before = append(hT1.Execute.Before, taskOriginT1)
	hT1.Execute.After = append(hT1.Execute.After, taskOriginT2)
	t1.Hooks = hT1

	return t1
}

func getTaskOther(name string) *Task {
	other := &Task{
		Name:     name,
		Playbook: "Playbook_overwritten",
	}
	other.EnvVars = make(map[string]string)
	other.EnvVars["key2"] = "val2_other"
	other.EnvVars["key3"] = "val3_other"
	other.Parameters = make(map[string]interface{})
	other.Parameters["key2"] = "val2_other"
	other.Parameters["key3"] = "val3_other"

	other.cRef = componentRef{
		ref: "cOther",
	}

	taskOtherT1 = TaskRef{ref: "T3"}
	taskOtherT2 = TaskRef{ref: "T4"}
	oT1 := TaskHook{}
	oT1.Execute.Before = append(oT1.Execute.Before, taskOtherT1)
	oT1.Execute.After = append(oT1.Execute.After, taskOtherT2)
	other.Hooks = oT1

	return other
}

func checkTaskMerge(t *testing.T, ta *Task) {

	assert.Equal(t, ta.cRef.ref, "cOther")
	assert.Equal(t, ta.Playbook, "Playbook_overwritten")

	if assert.Len(t, ta.EnvVars, 3) {
		checkMap(t, ta.EnvVars, "key1", "val1_target")
		checkMap(t, ta.EnvVars, "key2", "val2_other")
		checkMap(t, ta.EnvVars, "key3", "val3_other")
	}

	if assert.Len(t, ta.Parameters, 3) {
		checkMapInterface(t, ta.Parameters, "key1", "val1_target")
		checkMapInterface(t, ta.Parameters, "key2", "val2_other")
		checkMapInterface(t, ta.Parameters, "key3", "val3_other")
	}

	if assert.Len(t, ta.Hooks.Execute.Before, 2) {
		assert.Contains(t, ta.Hooks.Execute.Before, taskOriginT1, taskOtherT1)
		assert.Equal(t, ta.Hooks.Execute.Before[0], taskOriginT1)
		assert.Equal(t, ta.Hooks.Execute.Before[1], taskOtherT1)
	}

	if assert.Len(t, ta.Hooks.Execute.After, 2) {
		assert.Contains(t, ta.Hooks.Execute.After, taskOriginT2, taskOtherT2)
		assert.Equal(t, ta.Hooks.Execute.After[0], taskOriginT2)
		assert.Equal(t, ta.Hooks.Execute.After[1], taskOtherT2)
	}

}

func TestTaskMerge(t *testing.T) {
	o := getTaskOrigin()
	err := o.customize(*getTaskOther("my_name"))
	if assert.Nil(t, err) {
		checkTaskMerge(t, o)
	}
}

func TestMergeTaskItself(t *testing.T) {
	o := getTaskOrigin()
	oi := o
	err := o.customize(*o)
	if assert.Nil(t, err) {
		assert.Equal(t, oi, o)
	}
}

func TestTasksMerge(t *testing.T) {
	origins := make(Tasks)
	origins["myS"] = getTaskOrigin()
	others := make(Tasks)
	others["myS"] = getTaskOther("my_name")

	customized, err := origins.customize(&Environment{}, others)
	if assert.Nil(t, err) {
		if assert.Len(t, customized, 1) {
			o := customized["myS"]
			checkTaskMerge(t, o)
		}
	}
}

func TestTasksMergeAddition(t *testing.T) {
	origins := make(Tasks)
	origins["myS"] = getTaskOrigin()
	others := make(Tasks)
	others["myS"] = getTaskOther("my_name")
	others["new"] = getTaskOther("new")

	customized, err := origins.customize(&Environment{}, others)
	if assert.Nil(t, err) {
		assert.Len(t, customized, 2)
	}
}

func TestTasksEmptyMerge(t *testing.T) {
	origins := make(Tasks)
	o := getTaskOrigin()
	origins["myS"] = o
	others := make(Tasks)

	customized, err := origins.customize(&Environment{}, others)
	if assert.Nil(t, err) {
		if assert.Len(t, customized, 1) {
			oc := customized["myS"]
			assert.Equal(t, o, oc)
		}
	}
}

func TestMergeTaskUnrelated(t *testing.T) {
	ta := Task{
		Name:       "Name",
		Parameters: Parameters{},
	}
	ta.Parameters["p1"] = "val1"

	o := Task{
		Name:       "Dummy",
		Parameters: Parameters{},
	}
	o.Parameters["p1"] = "val1"

	err := ta.customize(o)
	if assert.NotNil(t, err) {
		assert.Equal(t, err.Error(), "cannot customize unrelated tasks (Name != Dummy)")
	}
	assert.Equal(t, 1, len(ta.Parameters))
	checkMapInterface(t, ta.Parameters, "p1", "val1")

}
