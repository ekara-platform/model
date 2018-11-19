package model

import (
	"encoding/json"
	"errors"
)

type NodeHook struct {
	Provision Hook
	Destroy   Hook
}

func (r NodeHook) MarshalJSON() ([]byte, error) {
	t := struct {
		Provision *Hook `json:",omitempty"`
		Destroy   *Hook `json:",omitempty"`
	}{}
	if r.Provision.HasTasks() {
		t.Provision = &r.Provision
	}
	if r.Destroy.HasTasks() {
		t.Destroy = &r.Destroy
	}
	return json.Marshal(t)
}

func (r NodeHook) HasTasks() bool {
	return r.Provision.HasTasks() ||
		r.Destroy.HasTasks()
}

func (r NodeHook) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	vErrs.merge(r.Provision.validate())
	vErrs.merge(r.Destroy.validate())
	return vErrs
}

func (r *NodeHook) merge(other NodeHook) error {
	if err := r.Provision.merge(other.Provision); err != nil {
		return err
	}
	if err := r.Destroy.merge(other.Destroy); err != nil {
		return err
	}
	return nil
}

// NodeSet contains the whole specification of a nodes et to create on a specific
// cloud provider
type NodeSet struct {
	// The name of the machines
	Name string
	// The number of machines to create
	Instances int
	// The ref to the provider where to create the machines
	Provider ProviderRef
	// The parameters related to the orchestrator used to manage the machines
	Orchestrator OrchestratorRef
	// Volumes attached to each node
	Volumes Volumes
	// Hooks for executing tasks around provisioning and destruction
	Hooks NodeHook
	// The labels associated with the nodeset
	Labels Labels
}

func (r NodeSet) DescType() string {
	return "NodeSet"
}

func (r NodeSet) DescName() string {
	return r.Name
}

func (r NodeSet) MarshalJSON() ([]byte, error) {
	provider, e := r.Provider.Resolve()
	if e != nil {
		return nil, e
	}
	orchestrator, e := r.Orchestrator.Resolve()
	if e != nil {
		return nil, e
	}
	t := struct {
		Name         string       `json:",omitempty"`
		Instances    int          `json:",omitempty"`
		Provider     Provider     `json:",omitempty"`
		Orchestrator Orchestrator `json:",omitempty"`
		Volumes      Volumes
		Hooks        *NodeHook `json:",omitempty"`
	}{
		Name:         r.Name,
		Instances:    r.Instances,
		Provider:     provider,
		Orchestrator: orchestrator,
		Volumes:      r.Volumes,
	}
	if r.Hooks.HasTasks() {
		t.Hooks = &r.Hooks
	}
	return json.Marshal(t)
}

func (r NodeSet) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	vErrs.merge(r.Provider.validate())
	vErrs.merge(r.Orchestrator.validate())
	vErrs.merge(r.Hooks.validate())
	return vErrs
}

func (r *NodeSet) merge(other NodeSet) error {
	if r.Name != other.Name {
		return errors.New("cannot merge unrelated node sets (" + r.Name + " != " + other.Name + ")")
	}
	if err := r.Provider.merge(other.Provider); err != nil {
		return err
	}
	if err := r.Orchestrator.merge(other.Orchestrator); err != nil {
		return err
	}
	if err := r.Volumes.merge(other.Volumes); err != nil {
		return err
	}
	if err := r.Hooks.merge(other.Hooks); err != nil {
		return err
	}
	if r.Instances < other.Instances {
		r.Instances = other.Instances
	}
	r.Labels = r.Labels.inherits(other.Labels)
	return nil
}

type NodeSets map[string]NodeSet

func createNodeSets(env *Environment, yamlEnv *yamlEnvironment) NodeSets {
	res := NodeSets{}
	for name, yamlNodeSet := range yamlEnv.Nodes {
		if yamlNodeSet.Instances <= 0 {
			env.errors.addError(errors.New("instances must be a positive number"), env.location.appendPath("nodes."+name+".instances"))
		}
		res[name] = NodeSet{
			Name:         name,
			Instances:    yamlNodeSet.Instances,
			Provider:     createProviderRef(env, env.location.appendPath("nodes."+name+".provider.name"), yamlNodeSet.Provider),
			Orchestrator: createOrchestratorRef(env, env.location.appendPath("nodes."+name+".orchestrator"), yamlNodeSet.Orchestrator),
			Volumes:      createVolumes(env, env.location.appendPath("nodes."+name+".volumes"), yamlNodeSet.Volumes),
			Hooks: NodeHook{
				Provision: createHook(env, env.location.appendPath("nodes."+name+".hooks.provision"), yamlNodeSet.Hooks.Provision),
				Destroy:   createHook(env, env.location.appendPath("nodes."+name+".hooks.destroy"), yamlNodeSet.Hooks.Destroy),
			},
			Labels: yamlNodeSet.Labels,
		}
	}
	return res
}

func (r NodeSets) merge(env *Environment, other NodeSets) error {
	for id, n := range other {
		if nodeSet, ok := r[id]; ok {
			if err := nodeSet.merge(n); err != nil {
				return err
			}
		} else {
			n.Provider.env = env
			n.Orchestrator.env = env
			r[id] = n
		}
	}
	return nil
}
