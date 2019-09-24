package model

import (
	"errors"
)

const (
	//GenericNodeSetName is the name of the generic node set
	//
	//The generic node set is intended to be used for sharing common
	// content, example: parameter, environment variables..., with all
	// others node sets within the whole descriptor.
	GenericNodeSetName = "*"
)

type (
	// NodeSet contains the whole specification of a nodeset to create on a specific
	// cloud provider
	NodeSet struct {
		location DescriptorLocation
		// The name of the machines
		Name string
		// The number of machines to create
		Instances int
		// The ref to the provider where to create the machines
		Provider ProviderRef
		// The parameters related to the orchestrator used to manage the machines
		Orchestrator OrchestratorRef
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

func (r NodeSet) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	if r.Instances <= 0 {
		vErrs.addError(errors.New("instances must be a positive number"), r.location.appendPath("instances"))
	}
	vErrs.merge(ErrorOnInvalid(r.Provider, r.Orchestrator, r.Hooks))
	return vErrs
}

func (r *NodeSet) customize(with NodeSet) error {
	if r.Name != with.Name && with.Name != GenericNodeSetName {
		return errors.New("cannot customize unrelated node sets (" + r.Name + " != " + with.Name + ")")
	}
	if err := r.Provider.customize(with.Provider); err != nil {
		return err
	}
	if err := r.Orchestrator.customize(with.Orchestrator); err != nil {
		return err
	}
	if err := r.Hooks.customize(with.Hooks); err != nil {
		return err
	}
	if with.Instances > 0 {
		r.Instances = with.Instances
	}
	r.Labels = r.Labels.inherit(with.Labels)
	return nil
}

func createNodeSets(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) (NodeSets, error) {
	// we will keep a reference on an eventual generic node set
	var gNs *NodeSet
	var err error
	res := NodeSets{}
	for name, yamlNodeSet := range yamlEnv.Nodes {
		nodeSetLocation := location.appendPath(name)
		if name == GenericNodeSetName {
			//The generic node set has been located
			gNs, err = buildNode(name, env, nodeSetLocation, yamlNodeSet)
			if err != nil {
				return NodeSets{}, err
			}
		} else {
			nodeSet, err := buildNode(name, env, nodeSetLocation, yamlNodeSet)
			if err != nil {
				return NodeSets{}, err
			}
			res[name] = *nodeSet
		}
	}

	if gNs != nil {
		// The generic node set will be merged into all others
		// in order to propagate the common stuff.
		for name, n := range res {
			err := n.customize(*gNs)
			if err != nil {
				return res, err
			}
			res[name] = n
		}
	}
	return res, nil
}

func buildNode(name string, env *Environment, location DescriptorLocation, yN yamlNode) (*NodeSet, error) {
	pRef, err := createProviderRef(env, location.appendPath("provider"), yN.Provider)
	if err != nil {
		return nil, err
	}
	oRef, err := createOrchestratorRef(env, location.appendPath("orchestrator"), yN.Orchestrator)
	if err != nil {
		return nil, err
	}
	pHook, err := createHook(env, location.appendPath("hooks.provision"), yN.Hooks.Provision)
	if err != nil {
		return nil, err
	}
	return &NodeSet{
		location:     location,
		Name:         name,
		Instances:    yN.Instances,
		Provider:     pRef,
		Orchestrator: oRef,
		Hooks: NodeHook{
			Provision: pHook,
		},
		Labels: yN.Labels}, nil
}

func (r NodeSets) customize(env *Environment, with NodeSets) (NodeSets, error) {
	res := make(map[string]NodeSet)
	for kr, vr := range r {
		res[kr] = vr
	}

	for id, n := range with {
		if nodeSet, ok := res[id]; ok {
			nm := &nodeSet
			if err := nm.customize(n); err != nil {
				return res, err
			}
			res[id] = *nm
		} else {
			n.Provider.env = env
			n.Orchestrator.env = env
			res[id] = n
		}
	}
	return res, nil
}
