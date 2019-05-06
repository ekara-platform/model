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
	GenericNodeSetName = "_"
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

func (r NodeSet) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	if r.Instances <= 0 {
		vErrs.addError(errors.New("instances must be a positive number"), r.location.appendPath("instances"))
	}
	vErrs.merge(ErrorOnInvalid(r.Provider, r.Orchestrator, r.Hooks, r.Volumes))
	return vErrs
}

func (r *NodeSet) merge(other NodeSet) error {
	if r.Name != other.Name && other.Name != GenericNodeSetName {
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

func createNodeSets(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) (NodeSets, error) {
	// we will keep a reference on an eventual generic node set
	var gNs *NodeSet
	res := NodeSets{}
	for name, yamlNodeSet := range yamlEnv.Nodes {
		nodeSetLocation := location.appendPath(name)
		if name == GenericNodeSetName {
			//The generic node set has been located
			gNs = buildNode(name, env, nodeSetLocation, yamlNodeSet)
		} else {
			res[name] = *buildNode(name, env, nodeSetLocation, yamlNodeSet)
		}
	}

	if gNs != nil {
		// The generic node set will be merged into all others
		// in order to propagate the common stuff.
		for name, n := range res {
			err := n.merge(*gNs)
			if err != nil {
				return res, err
			}
			res[name] = n
		}
	}

	for _, n := range res {
		env.Ekara.tagUsedComponent(n.Provider)
	}

	return res, nil
}

func buildNode(name string, env *Environment, location DescriptorLocation, yN yamlNode) *NodeSet {
	return &NodeSet{
		location:     location,
		Name:         name,
		Instances:    yN.Instances,
		Provider:     createProviderRef(env, location.appendPath("provider"), yN.Provider),
		Orchestrator: createOrchestratorRef(env, location.appendPath("orchestrator"), yN.Orchestrator),
		Volumes:      createVolumes(location.appendPath("volumes"), yN.Volumes),
		Hooks: NodeHook{
			Provision: createHook(env, location.appendPath("hooks.provision"), yN.Hooks.Provision),
			Destroy:   createHook(env, location.appendPath("hooks.destroy"), yN.Hooks.Destroy),
		},
		Labels: yN.Labels,
	}
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
