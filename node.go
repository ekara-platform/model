package descriptor

import (
	"errors"
	"strings"
)

type NodeSet struct {
	root *Environment
	Labels

	Name      string
	Provider  ProviderRef
	Instances int

	Hooks struct {
		Provision Hook
		Destroy   Hook
	}
}

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
			nodeSet := NodeSet{
				root:      env,
				Labels:    createLabels(vErrs, yamlNodeSet.Labels...),
				Name:      name,
				Instances: yamlNodeSet.Instances}

			nodeSet.Provider = createProviderRef(vErrs, env, "nodes."+name+".provider", yamlNodeSet.Provider)
			nodeSet.Hooks.Provision = createHook(vErrs, env.Tasks, "nodes."+name+".hooks.provision", yamlNodeSet.Hooks.Provision)
			nodeSet.Hooks.Destroy = createHook(vErrs, env.Tasks, "nodes."+name+".hooks.destroy", yamlNodeSet.Hooks.Destroy)

			res[name] = nodeSet
		}
	}
	return res
}

func createNodeSetRef(vErrs *ValidationErrors, env *Environment, location string, labels ... string) NodeSetRef {
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
