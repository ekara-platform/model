package model

import (
	"log"

	"net/url"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

// yaml tag for labels
type yamlLabels struct {
	Labels []string
}

// yaml tag for parameters
type yamlParams struct {
	Params attributes
}

// yaml tag for docker parameters
type yamlDocker struct {
	Docker attributes
}

// yaml reference on a name allowing to hold more specific parameters
type yamlRef struct {
	yamlParams `yaml:",inline"`
	Name       string
}

// yaml reference on a provider name allowing to hold more specific parameters
type yamlProviderRef struct {
	// The name of the referenced provider
	Name string
	// The overwritten parameters for the provider
	yamlParams `yaml:",inline"`
	// The volumes to create and mount
	Volumes []yamlVolumes
}

//yaml tag for a volume and its parameters
type yamlVolumes struct {
	// The mounting path of the created volume
	Name string
	// The parameters required to create the volume.
	// These parameters are typically provider dependent, so refer to the provider documentation to figure how to create volumes.
	yamlParams `yaml:",inline"`
}

//yaml tag for hooks
type yamlHook struct {
	// Hooks to be executed before the corresponding process step
	Before []yamlRef
	// Hooks to be executed after the corresponding process step
	After []yamlRef
}

// yaml tag for a repository and its version
type yamlRepoVersion struct {
	// The Repository, following the notation "organization/repo_name"
	Repository string
	// The release into the repository
	Version string
}

// Definition of the Lagoon environment
type yamlEnvironment struct {
	// The name of the environment
	Name string
	// The description of the environment
	Description string
	// The version of the environment
	Version string

	// The labels associated to the environment
	yamlLabels `yaml:",inline"`

	// Settings
	Settings struct {
		ComponentBase  string `yaml:"componentBase"`
		DockerRegistry string `yaml:"dockerRegistry"`
		Proxy          struct {
			Http    string
			Https   string
			NoProxy string `yaml:"noProxy"`
		}
	}

	// The Lagoon platform used to interact with the environment
	LagoonPlatform yamlRepoVersion `yaml:"lagoonPlatform"`

	// Imports, to be included into the environment descriptor
	Imports []string

	// Components
	Components map[string]string

	// Global definition of the orchestrator to install on the environment
	Orchestrator struct {
		// The orchestrator specifics parameters
		yamlParams `yaml:",inline"`
		// The docker parameters
		yamlDocker `yaml:",inline"`
		// The name of the orchestrator
		Name string
		// The repository and version of the orchestrator
		yamlRepoVersion `yaml:",inline"`
	}

	// The list of all cloud providers required to create the environment
	Providers map[string]struct {
		// The provider parameters
		yamlParams `yaml:",inline"`
		// The repository and version of the provider
		yamlRepoVersion `yaml:",inline"`
	}

	// The list of node sets to create
	Nodes map[string]struct {
		// The labels associated to the node set
		yamlLabels `yaml:",inline"`
		// Reference on the provider where to create the node set
		Provider yamlProviderRef
		// The number of instances to create within the node set
		Instances int

		// The installed orchestrator
		Orchestrator struct {
			// The overwritten orchestrator specifics parameters for the node set
			yamlParams `yaml:",inline"`
			// The overwritten docker parameters for the node set
			yamlDocker `yaml:",inline"`
		}

		// The Hooks to be executed while provisionning and destoying the node set
		Hooks struct {
			Provision yamlHook
			Destroy   yamlHook
		}
	}

	// Software stacks to be installed on the environment
	Stacks map[string]struct {
		// The labels associated to the stack
		yamlLabels `yaml:",inline"`
		// The repository and version of the stack
		yamlRepoVersion `yaml:",inline"`
		// The names of the node sets where the stack must de installed
		DeployOn []string `yaml:"deployOn"`

		// The Hooks to be executed while deploying and undeploying the stack
		Hooks struct {
			Deploy   yamlHook
			Undeploy yamlHook
		}
	}

	// Custom tasks which can be run on the created environment
	Tasks map[string]struct {
		// The labels associated to the task
		yamlLabels `yaml:",inline"`
		// The task parameters
		yamlParams `yaml:",inline"`

		// The name of the playbook to launch the task
		Playbook string
		// The CRON to run cyclically the task
		Cron string
		// The labels allowing to locate where to run the task ( on which node sets )
		RunOn []string `yaml:"runOn"`
		// The Hooks to be executed in addition the the main task playbook
		Hooks struct {
			Execute yamlHook
		}
	}

	// Global hooks
	Hooks struct {
		Init      yamlHook
		Provision yamlHook
		Deploy    yamlHook
		Undeploy  yamlHook
		Destroy   yamlHook
	}
}

func parseYamlDescriptor(logger *log.Logger, u *url.URL) (env yamlEnvironment, err error) {
	var normalizedUrl *url.URL
	normalizedUrl, err = NormalizeUrl(u)
	if err != nil {
		return
	}
	baseLocation, content, err := ReadUrl(logger, normalizedUrl)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(content, &env)
	if err != nil {
		return
	}

	err = processYamlImports(logger, baseLocation, &env)
	if err != nil {
		return
	}

	return
}

func processYamlImports(logger *log.Logger, base *url.URL, env *yamlEnvironment) error {
	if len(env.Imports) > 0 {
		for _, val := range env.Imports {
			importUrl, err := url.Parse(val)
			if err != nil {
				return err
			}
			ref := base.ResolveReference(importUrl)
			logger.Println("Processing import", ref)
			importedDesc, err := parseYamlDescriptor(logger, ref)
			if err != nil {
				return err
			}
			mergo.Merge(env, importedDesc)
		}
		env.Imports = nil
	} else {
		logger.Println("No import to process")
	}
	return nil
}
