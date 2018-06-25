package model

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v2"
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

	// TODO Document this
	Hooks struct {
		Provision Hook
		Destroy   Hook
	}
}

type NodeSetRef struct {
	nodeSets []*NodeSet
}

type Client struct {
	Name string
	Uid  string
}

type MachineConfig struct {
	ConnectionConfig ConnectionConfig `yaml:"connectionConfig"`
}

type ConnectionConfig struct {
	Provider          string
	MachinePublicKey  string `yaml:"machine_public_key"`
	MachinePrivateKey string `yaml:"machine_private_key"`
}

type NodeParams struct {
	Client           Client
	Params           map[string]interface{}
	Instances        int
	ConnectionConfig ConnectionConfig `yaml:"connectionConfig"`
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
//
//	Parameters:
//		c: the client requesting the creation
//		uid: the unique Id used to tag the created machines
//		p: the name of the provider where to create the machines
//		pubK: the SSH public key to connect the machine
//		privK: the SSH private key to connect the machine
func (n NodeSet) NodeParams(c string, uid string, p string, pubK string, privK string) (b []byte, e error) {
	cli := Client{Name: c, Uid: uid}
	nev := NodeParams{
		Client:    cli,
		Params:    n.Provider.Parameters.copy(),
		Instances: n.Instances,
	}

	mConf := ConnectionConfig{
		Provider:          p,
		MachinePublicKey:  pubK,
		MachinePrivateKey: privK,
	}
	nev.ConnectionConfig = mConf
	b, e = yaml.Marshal(&nev)
	return
}

// OrchestratorParams returns the parameters required to deploy the orchestrator
// on a nodeset
func (n NodeSet) OrchestratorParams() (b []byte, e error) {
	b, e = yaml.Marshal(n.Orchestrator.Parameters.copy())
	return
}
