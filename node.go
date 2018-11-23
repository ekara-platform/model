package model

import (
	"encoding/json"
	"errors"
)

type (
	// NodeSet contains the whole specification of a nodeset to create on a specific
	// cloud provider
	NodeSet struct {
		// The name of the machines
		Name string
		// The number of machines to create
		Instances int
		// The ref to the provider where to create the machines
		Provider providerRef
		// The parameters related to the orchestrator used to manage the machines
		Orchestrator orchestratorRef
		// Volumes attached to each node
		Volumes Volumes
		// The hooks linked to the node set lifecycle events
		Hooks NodeHook
		// The labels associated with the nodeset
		Labels Labels
	}

	//NodeSets represents all the node sets of the environment
	NodeSets map[string]NodeSet
)

//DescType returns the Describable type of the node set
//  Hardcoded to : "NodeSet"
func (r NodeSet) DescType() string {
	return "NodeSet"
}

//DescName returns the Describable name of the node set
func (r NodeSet) DescName() string {
	return r.Name
}

// MarshalJSON returns the serialized content of node set as JSON
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
	return ErrorOnInvalid(r.Provider, r.Orchestrator, r.Hooks)
}

func (r *NodeSet) merge(other NodeSet) error {
	if r.Name != other.Name {
		return errors.New("cannot merge unrelated node sets (" + r.Name + " != " + other.Name + ")")
	}
	if err := r.Provider.merge(other.Provider); err != nil {
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
