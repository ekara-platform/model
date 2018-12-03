package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
func (r *Task) merge(other Task) error {
	if r.Name != other.Name {
		return errors.New("cannot merge unrelated stacks (" + r.Name + " != " + other.Name + ")")
	}
	if err := r.Component.merge(other.Component); err != nil {
		return err
	}
	if err := r.Hooks.merge(other.Hooks); err != nil {
		return err
	}
	if r.Playbook == "" {
		r.Playbook = other.Playbook
	}
	if r.Cron == "" {
		r.Cron = other.Cron
	}
	r.Parameters = r.Parameters.inherits(other.Parameters)
	r.EnvVars = r.EnvVars.inherits(other.EnvVars)
	return nil
}
*/

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
		Parameters: Parameters{},
	}
	ta.Parameters["p1"] = "val1"

	err := ta.merge(ta)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(ta.Parameters))
	assert.Contains(t, ta.Parameters, "p1")
	assert.Equal(t, ta.Parameters["p1"], "val1")
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
	// Hooks have specific tests
	// Componebt has specific tests

	o := Task{
		Name:       "Name",
		Playbook:   "Playbook_updated",
		Cron:       "Cron_updated",
		Parameters: Parameters{},
		EnvVars:    EnvVars{},
	}
	o.Parameters["p1"] = "val1_updated"
	o.EnvVars["e1"] = "env1_updated"

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
}

func TestMergeTaskAddition(t *testing.T) {
	ta := Task{
		Name: "Name",

		Parameters: Parameters{},
		EnvVars:    EnvVars{},
	}
	ta.Parameters["p1"] = "val1"
	ta.EnvVars["e1"] = "env1"
	// Hooks have specific tests
	// Componebt has specific tests

	o := Task{
		Name:       "Name",
		Playbook:   "Playbook",
		Cron:       "Cron",
		Parameters: Parameters{},
		EnvVars:    EnvVars{},
	}
	o.Parameters["p2"] = "val2_added"
	o.EnvVars["e2"] = "env2_added"

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
	ts.merge(env, os)

	assert.Equal(t, 3, len(ts))

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
