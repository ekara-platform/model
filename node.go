package model

import (
	"encoding/json"
	"errors"
	"fmt"

	_ "gopkg.in/yaml.v2"
)

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
	Volumes []Volume
	// Hooks for executing tasks around provisioning and destruction
	Hooks NodeHook
	// The labels associated with the nodeset
	Labels Labels
}

func (n NodeSet) HumanDescribe() string {
	return fmt.Sprintf("NodeSet: %s", n.Name)
}

func (r NodeSet) MarshalJSON() ([]byte, error) {
	t := struct {
		Name         string       `json:",omitempty"`
		Instances    int          `json:",omitempty"`
		Provider     Provider     `json:",omitempty"`
		Orchestrator Orchestrator `json:",omitempty"`
		Volumes      []Volume
		Hooks        *NodeHook `json:",omitempty"`
	}{
		Name:         r.Name,
		Instances:    r.Instances,
		Provider:     r.Provider.Resolve(),
		Orchestrator: r.Orchestrator.Resolve(),
		Volumes:      r.Volumes,
	}
	if r.Hooks.HasTasks() {
		t.Hooks = &r.Hooks
	}
	return json.Marshal(t)
}

type NodeHook struct {
	Provision Hook
	Destroy   Hook
}

func (r NodeHook) HasTasks() bool {
	return r.Provision.HasTasks() ||
		r.Destroy.HasTasks()
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

// Reference to a node set
type NodeSetRef struct {
	nodeSets []*NodeSet
}

func createNodeSets(vErrs *ValidationErrors, env *Environment, yamlEnv *yamlEnvironment) map[string]NodeSet {
	res := map[string]NodeSet{}
	if yamlEnv.Nodes == nil || len(yamlEnv.Nodes) == 0 {
		vErrs.AddError(errors.New("no node specified"), "nodes")
	} else {
		for name, yamlNodeSet := range yamlEnv.Nodes {
			if yamlNodeSet.Instances <= 0 {
				vErrs.AddError(errors.New("node set instances must be a positive number"), "nodes."+name+".instances")
			}

			res[name] = NodeSet{
				Name:         name,
				Instances:    yamlNodeSet.Instances,
				Provider:     createProviderRef(vErrs, "nodes."+name+".provider", env, yamlNodeSet.Provider),
				Orchestrator: createOrchestratorRef(env, yamlNodeSet.Orchestrator),
				Volumes:      createVolumes(vErrs, "nodes."+name+".volumes", yamlNodeSet.Volumes),
				Hooks: NodeHook{
					Provision: createHook(vErrs, "nodes."+name+".hooks.provision", env, yamlNodeSet.Hooks.Provision),
					Destroy:   createHook(vErrs, "nodes."+name+".hooks.destroy", env, yamlNodeSet.Hooks.Destroy),
				},
				Labels: yamlNodeSet.Labels,
			}
		}
	}
	return res
}

func createNodeSetRef(vErrs *ValidationErrors, env *Environment, location string, nodeSetRefs ...string) NodeSetRef {
	nodeSets := make([]*NodeSet, 0, 10)
	if len(nodeSetRefs) == 0 {
		for _, nodeSet := range env.NodeSets {
			nodeSets = append(nodeSets, &nodeSet)
		}
	} else {
		for _, nodeSetRef := range nodeSetRefs {
			if nodeSet, ok := env.NodeSets[nodeSetRef]; ok {
				nodeSets = append(nodeSets, &nodeSet)
			} else {
				vErrs.AddError(errors.New("unknown node set reference: "+nodeSetRef), location)
			}
		}
	}
	return NodeSetRef{nodeSets: nodeSets}
}
