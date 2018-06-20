package model

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v2"
)

type OrchestratorParameters struct {
	Parameters attributes
}

type NodeSet struct {
	root *Environment
	Labels

	Name         string
	Provider     ProviderRef
	Instances    int
	Orchestrator OrchestratorParameters

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
				root:         env,
				Labels:       createLabels(vErrs, yamlNodeSet.Labels...),
				Name:         name,
				Instances:    yamlNodeSet.Instances,
				Orchestrator: OrchestratorParameters{},
			}

			nodeSet.Orchestrator.Parameters = createAttributes(yamlNodeSet.Orchestrator.Params, env.Orchestrator.Parameters.copy())

			nodeSet.Provider = createProviderRef(vErrs, env, "nodes."+name+".provider", yamlNodeSet.Provider)
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

type Client struct {
	Name string
	Uid  string
}

type NodeConfig struct {
	Client            Client
	Params            map[string]interface{}
	Instances         int
	Provider          string
	MachinePublicKey  string `yaml:"machine_public_key"`
	MachinePrivateKey string `yaml:"machine_private_key"`
}

func (n NodeSet) Config(c string, uid string, p string, pubK string, privK string) (b []byte, e error) {
	cli := Client{Name: c, Uid: uid}
	nev := NodeConfig{
		Client:            cli,
		Params:            n.Provider.Parameters.copy(),
		Instances:         n.Instances,
		Provider:          p,
		MachinePrivateKey: privK,
		MachinePublicKey:  pubK,
	}
	b, e = yaml.Marshal(&nev)
	return
}

func (n NodeSet) OrchestratorVars() (b []byte, e error) {
	b, e = yaml.Marshal(n.Orchestrator.Parameters.copy())
	return
}
