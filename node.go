package model

import (
	"errors"
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
	Hooks struct {
		Provision Hook
		Destroy   Hook
	}
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
				Hooks: struct {
					Provision Hook
					Destroy   Hook
				}{
					Provision: createHook(vErrs, "nodes."+name+".hooks.provision", env, yamlNodeSet.Hooks.Provision),
					Destroy:   createHook(vErrs, "nodes."+name+".hooks.destroy", env, yamlNodeSet.Hooks.Destroy)}}
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
