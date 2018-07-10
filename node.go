package model

import (
	"errors"
	"strings"

	_ "gopkg.in/yaml.v2"
)

// OrchestratorParameters contains the parameters related to the orchestrator
// defined into the environment
type OrchestratorParameters struct {
	//The Orchestrator specific parameters
	Parameters attributes
	// The Dockers specific parameters
	Docker attributes
}

// NodeSet contains the whole specification of a Nodeset to create on a specific
// cloud provider
type NodeSet struct {
	// The environment holding the nodeset
	root *Environment
	// The severals labels used to tag the machines
	Labels

	// The name of the machines
	Name string
	// The ref to the provider where to create the machines
	Provider ProviderRef
	// The number of machines to create
	Instances int

	// The parameters related to the orchestrator used to manage the machines
	Orchestrator OrchestratorParameters

	Volumes []Volume

	// TODO Document this
	Hooks struct {
		Provision Hook
		Destroy   Hook
	}
}

type NodeSetRef struct {
	nodeSets []*NodeSet
}

type NodeParams struct {
	Params    map[string]interface{}
	Instances int
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

			nodeSet := NodeSet{
				root:         env,
				Labels:       createLabels(vErrs, yamlNodeSet.Labels...),
				Name:         name,
				Instances:    yamlNodeSet.Instances,
				Orchestrator: OrchestratorParameters{},
			}

			nodeSet.Orchestrator.Parameters = createAttributes(yamlNodeSet.Orchestrator.Params, env.Orchestrator.Parameters.copy())
			nodeSet.Orchestrator.Docker = createAttributes(yamlNodeSet.Orchestrator.Docker, env.Orchestrator.Docker.copy())

			nodeSet.Provider = createProviderRef(vErrs, env, "nodes."+name+".provider", yamlNodeSet.Provider)
			nodeSet.Volumes = createVolumes(vErrs, env, "nodes."+name+".volumes", yamlNodeSet.Volumes)
			nodeSet.Hooks.Provision = createHook(vErrs, env.Tasks, "nodes."+name+".hooks.provision", yamlNodeSet.Hooks.Provision)
			nodeSet.Hooks.Destroy = createHook(vErrs, env.Tasks, "nodes."+name+".hooks.destroy", yamlNodeSet.Hooks.Destroy)

			res[name] = nodeSet
		}
	}
	return res
}

func createNodeSetRef(vErrs *ValidationErrors, env *Environment, location string, labels ...string) NodeSetRef {
	nodeSets := make([]*NodeSet, 0, 10)
	if len(labels) == 0 {
		vErrs.AddError(errors.New("empty node set reference"), location)
	} else {
		for _, nodeSet := range env.NodeSets {
			if nodeSet.MatchesLabels(labels...) {
				nodeSets = append(nodeSets, &nodeSet)
			}
		}
		if len(nodeSets) == 0 {
			vErrs.AddError(errors.New("no node set matches label(s): "+strings.Join(labels, ", ")), location)
		}
	}
	return NodeSetRef{nodeSets: nodeSets}
}

// NodeParams returns the parameters required to create a nodeset
func (n NodeSet) NodeParams() NodeParams {
	r := NodeParams{
		Params:    n.Provider.Parameters.copy(),
		Instances: n.Instances,
	}
	return r
}

// OrchestratorParams returns the parameters required to install the orchestrator
func (n NodeSet) OrchestratorParams() map[string]interface{} {
	r := make(map[string]interface{})
	r["docker"] = n.Orchestrator.Docker.copy()
	r["params"] = n.Orchestrator.Parameters.copy()
	return r
}
