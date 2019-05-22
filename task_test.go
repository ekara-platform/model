package model

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	err := ta.merge(o)
	if assert.NotNil(t, err) {
		assert.Equal(t, err.Error(), "cannot merge unrelated tasks (Name != Dummy)")
	}
	assert.Equal(t, 1, len(ta.Parameters))
	assert.Contains(t, ta.Parameters, "p1")
	assert.Equal(t, ta.Parameters["p1"], "val1")
}

func TestMergeTaskItself(t *testing.T) {
	ta := Task{
		Name:       "Name",
		Playbook:   "Playbook",
		Cron:       "Cron",
		Parameters: Parameters{},
		EnvVars:    EnvVars{},
	}
	ta.Parameters["p1"] = "val1"
	ta.EnvVars["e1"] = "env1"

	ta.cRef = componentRef{
		ref:       "cRef",
		mandatory: true,
	}

	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := TaskHook{}
	h.Execute.Before = append(h.Execute.Before, task1)
	h.Execute.After = append(h.Execute.After, task2)

	ta.Hooks = h

	err := ta.merge(ta)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(ta.Parameters))
	assert.Contains(t, ta.Parameters, "p1")
	assert.Equal(t, ta.Parameters["p1"], "val1")

	assert.Equal(t, 1, len(ta.EnvVars))
	assert.Contains(t, ta.EnvVars, "e1")
	assert.Equal(t, ta.EnvVars["e1"], "env1")

	assert.Equal(t, ta.cRef.ref, "cRef")
	assert.True(t, ta.cRef.mandatory)

	assert.True(t, reflect.DeepEqual(h, ta.Hooks))
}

func TestMergeTaskNoUpdate(t *testing.T) {
	ta := Task{
		Name:       "Name",
		Playbook:   "Playbook",
		Cron:       "Cron",
		Parameters: Parameters{},
		EnvVars:    EnvVars{},
	}
	ta.Parameters["p1"] = "val1"
	ta.EnvVars["e1"] = "env1"

	ta.cRef = componentRef{
		ref:       "cRef",
		mandatory: true,
	}

	task1 := TaskRef{ref: "ref1"}
	task2 := TaskRef{ref: "ref2"}
	h := TaskHook{}
	h.Execute.Before = append(h.Execute.Before, task1)
	h.Execute.After = append(h.Execute.After, task2)

	ta.Hooks = h

	o := Task{
		Name:       "Name",
		Playbook:   "Playbook_updated",
		Cron:       "Cron_updated",
		Parameters: Parameters{},
		EnvVars:    EnvVars{},
	}
	o.Parameters["p1"] = "val1_updated"
	o.EnvVars["e1"] = "env1_updated"

	o.cRef = componentRef{
		ref:       "cRef_updated",
		mandatory: false,
	}

	tasko1 := TaskRef{ref: "ref1_updated"}
	tasko2 := TaskRef{ref: "ref2_updated"}
	ho := TaskHook{}
	ho.Execute.Before = append(ho.Execute.Before, tasko1)
	ho.Execute.After = append(ho.Execute.After, tasko2)

	o.Hooks = ho

	err := ta.merge(o)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(ta.Parameters))
	assert.Contains(t, ta.Parameters, "p1")
	assert.Equal(t, ta.Parameters["p1"], "val1")
	assert.Equal(t, 1, len(ta.EnvVars))
	assert.Contains(t, ta.EnvVars, "e1")
	assert.Equal(t, ta.EnvVars["e1"], "env1")
	assert.Equal(t, ta.Playbook, "Playbook")
	assert.Equal(t, ta.Cron, "Cron")

	// The component should not be updated
	assert.Equal(t, ta.cRef.ref, "cRef")
	assert.True(t, ta.cRef.mandatory)

	// The hook should be updated with the news tasks
	if assert.False(t, reflect.DeepEqual(h, ta.Hooks)) {
		assert.Equal(t, 2, len(ta.Hooks.Execute.Before))
		assert.Equal(t, 2, len(ta.Hooks.Execute.After))
	}
}

func TestMergeTaskAddition(t *testing.T) {
	ta := Task{
		Name: "Name",

		Parameters: Parameters{},
		EnvVars:    EnvVars{},
	}
	ta.Parameters["p1"] = "val1"
	ta.EnvVars["e1"] = "env1"

	o := Task{
		Name:       "Name",
		Playbook:   "Playbook",
		Cron:       "Cron",
		Parameters: Parameters{},
		EnvVars:    EnvVars{},
	}
	o.Parameters["p2"] = "val2_added"
	o.EnvVars["e2"] = "env2_added"

	o.cRef = componentRef{
		ref:       "cRef_added",
		mandatory: false,
	}

	tasko1 := TaskRef{ref: "ref1_added"}
	tasko2 := TaskRef{ref: "ref2_added"}
	ho := TaskHook{}
	ho.Execute.Before = append(ho.Execute.Before, tasko1)
	ho.Execute.After = append(ho.Execute.After, tasko2)

	o.Hooks = ho

	err := ta.merge(o)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(ta.Parameters))
	assert.Contains(t, ta.Parameters, "p1")
	assert.Equal(t, ta.Parameters["p1"], "val1")
	assert.Contains(t, ta.Parameters, "p2")
	assert.Equal(t, ta.Parameters["p2"], "val2_added")
	assert.Equal(t, 2, len(ta.EnvVars))
	assert.Contains(t, ta.EnvVars, "e1")
	assert.Equal(t, ta.EnvVars["e1"], "env1")
	assert.Contains(t, ta.EnvVars, "e2")
	assert.Equal(t, ta.EnvVars["e2"], "env2_added")
	assert.Equal(t, ta.Playbook, "Playbook")
	assert.Equal(t, ta.Cron, "Cron")

	// The component should be added
	assert.Equal(t, ta.cRef.ref, "cRef_added")
	assert.False(t, ta.cRef.mandatory)

	// The hook should be updated with the news tasks
	assert.Equal(t, 1, len(ta.Hooks.Execute.Before))
	assert.Equal(t, 1, len(ta.Hooks.Execute.After))
}

func TestMergeNoTasks(t *testing.T) {
	ta1 := &Task{
		Name: "Name1",

		Parameters: Parameters{},
		EnvVars:    EnvVars{},
	}
	ta1.Parameters["p11"] = "val11"
	ta1.EnvVars["e11"] = "env11"

	ta2 := &Task{
		Name: "Name2",

		Parameters: Parameters{},
		EnvVars:    EnvVars{},
	}
	ta2.Parameters["p12"] = "val12"
	ta2.EnvVars["e12"] = "env12"

	ts := Tasks{}
	ts[ta1.Name] = ta1
	ts[ta2.Name] = ta2

	emptyTs := Tasks{}

	env := &Environment{}
	ts.merge(env, emptyTs)
	assert.Equal(t, 2, len(ts))
}

func TestMergeTasks(t *testing.T) {
	ta1 := &Task{
		Name: "Name1",

		Parameters: Parameters{},
		EnvVars:    EnvVars{},
	}
	ta1.Parameters["p11"] = "val11"
	ta1.EnvVars["e11"] = "env11"

	ta2 := &Task{
		Name: "Name2",

		Parameters: Parameters{},
		EnvVars:    EnvVars{},
	}
	ta2.Parameters["p12"] = "val12"
	ta2.EnvVars["e12"] = "env12"

	ts := Tasks{}
	ts[ta1.Name] = ta1
	ts[ta2.Name] = ta2

	o1 := &Task{
		Name: "Name1",

		Parameters: Parameters{},
		EnvVars:    EnvVars{},
	}
	o1.Parameters["p11"] = "update" // Not supposed to be merge
	o1.Parameters["p12"] = "new"    // Must be merged
	o1.EnvVars["e11"] = "update"    // Not supposed to be merge
	o1.EnvVars["e12"] = "new"       // Must be merged

	o3 := &Task{ // The whole task is supposed to be merged
		Name: "Name3",

		Parameters: Parameters{},
		EnvVars:    EnvVars{},
	}
	o3.Parameters["p13"] = "val13"
	o3.EnvVars["e13"] = "env13"

	os := Tasks{}
	os[o1.Name] = o1
	os[o3.Name] = o3

	env := &Environment{}
	ts, err := ts.merge(env, os)
	assert.Nil(t, err)

	if assert.Equal(t, 3, len(ts)) {

		if assert.Equal(t, 2, len(ts[ta1.Name].Parameters)) {
			assert.Equal(t, "val11", ts[ta1.Name].Parameters["p11"])
			assert.Equal(t, "new", ts[ta1.Name].Parameters["p12"])
		}
		if assert.Equal(t, 2, len(ts[ta1.Name].EnvVars)) {
			assert.Equal(t, "env11", ts[ta1.Name].EnvVars["e11"])
			assert.Equal(t, "new", ts[ta1.Name].EnvVars["e12"])
		}

		if assert.Equal(t, 1, len(ts[ta2.Name].Parameters)) {
			assert.Equal(t, "val12", ts[ta2.Name].Parameters["p12"])
		}
		if assert.Equal(t, 1, len(ts[ta2.Name].EnvVars)) {
			assert.Equal(t, "env12", ts[ta2.Name].EnvVars["e12"])
		}

		if assert.Equal(t, 1, len(ts[o3.Name].Parameters)) {
			assert.Equal(t, "val13", ts[o3.Name].Parameters["p13"])

		}
		if assert.Equal(t, 1, len(ts[o3.Name].EnvVars)) {
			assert.Equal(t, "env13", ts[o3.Name].EnvVars["e13"])
		}
	}
}
