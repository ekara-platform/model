package model

import (
	"log"

	"net/url"

	"bytes"
	"text/template"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

// yaml tag for the proxy details
type yamlProxy struct {
	Http    string
	Https   string
	NoProxy string `yaml:"noProxy"`
}

// yaml tag for parameters
type yamlParams struct {
	Params map[string]interface{}
}

// yaml tag for Docker parameters
type yamlDockerParams struct {
	Docker map[string]interface{}
}

// yaml tag for environment variables
type yamlEnv struct {
	Env map[string]string
}

// yaml tag for component
type yamlComponent struct {
	Repository string
	Version    string
}

// yaml tag for a volume and its parameters
type yamlVolume struct {
	// The mounting path of the created volume
	Path string
	// The parameters required to create the volume (typically provider dependent)
	yamlParams `yaml:",inline"`
}

// yaml reference to provider
type yamlProviderRef struct {
	Name string
	// The overriding provider parameters
	yamlParams `yaml:",inline"`
	// The overriding provider environment variables
	yamlEnv `yaml:",inline"`
	// The overriding provider proxy
	Proxy yamlProxy
}

// yaml reference to orchestrator
type yamlOrchestratorRef struct {
	// The overriding orchestrator parameters
	yamlParams `yaml:",inline"`
	// The overriding docker parameters
	yamlDockerParams `yaml:",inline"`
	// The overriding orchestrator environment variables
	yamlEnv `yaml:",inline"`
}

// yaml reference to task
type yamlTaskRef struct {
	// The referenced task
	Task string
	// The overriding parameters
	yamlParams `yaml:",inline"`
	// The overriding environment variables
	yamlEnv `yaml:",inline"`
}

//yaml tag for hooks
type yamlHook struct {
	// Hooks to be executed before the corresponding process step
	Before []yamlTaskRef
	// Hooks to be executed after the corresponding process step
	After []yamlTaskRef
}

func (e *yamlEnvironment) RawContent() ([]byte, error) {
	return yaml.Marshal(e)
}

// Definition of the Lagoon environment
type yamlEnvironment struct {
	// Imports, to be included into the environment descriptor
	Imports []string

	// The name of the environment
	Name string
	// The description of the environment
	Description string

	Version string

	// The Lagoon platform used to interact with the environment
	Lagoon struct {
		ComponentBase  string `yaml:"componentBase"`
		DockerRegistry string `yaml:"dockerRegistry"`
		Components     map[string]yamlComponent
	}

	// Tasks which can be run on the created environment
	Tasks map[string]struct {
		// The task parameters
		yamlParams `yaml:",inline"`
		// The task environment variables
		yamlEnv `yaml:",inline"`
		// The name of the playbook to launch the task
		Playbook string
		// The name of the node sets to run the task on (all if not specified)
		On []string
		// The CRON to run cyclically the task
		Cron string
		// The Hooks to be executed in addition the the main task playbook
		Hooks struct {
			Execute yamlHook
		}
	}

	// Global definition of the orchestrator to install on the environment
	Orchestrator struct {
		// Name of the orchestrator component
		Component string
		// The orchestrator parameters
		yamlParams `yaml:",inline"`
		// The orchestrator environment variables
		yamlEnv `yaml:",inline"`
		// The Docker parameters
		yamlDockerParams `yaml:",inline"`
	}

	// The list of all cloud providers required to create the environment
	Providers map[string]struct {
		// Name of the provider component
		Component string
		// The provider parameters
		yamlParams `yaml:",inline"`
		// The provider environment variables
		yamlEnv `yaml:",inline"`
		// The provider proxy
		Proxy yamlProxy
	}

	// The list of node sets to create
	Nodes map[string]struct {
		// The number of instances to create within the node set
		Instances int
		// The provider used to create the node set and its settings
		Provider yamlProviderRef
		// The orchestrator settings for this node set
		Orchestrator yamlOrchestratorRef
		// The orchestrator settings for this node set
		Volumes []yamlVolume
		// The Hooks to be executed while provisioning and destroying the node set
		Hooks struct {
			Provision yamlHook
			Destroy   yamlHook
		}
	}

	// Software stacks to be installed on the environment
	Stacks map[string]struct {
		// Name of the stack component
		Component string
		// Name of the nodes to deploy the stack on (all if not specified)
		On []string
		// The Hooks to be executed while deploying and undeploying the stack
		Hooks struct {
			Deploy   yamlHook
			Undeploy yamlHook
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

func parseYamlDescriptor(logger *log.Logger, u *url.URL, data map[string]interface{}) (env yamlEnvironment, err error) {
	var normalizedUrl *url.URL
	normalizedUrl, err = NormalizeUrl(u)
	if err != nil {
		return
	}

	// Read descriptor content
	baseLocation, content, err := ReadUrl(logger, normalizedUrl)
	if err != nil {
		return
	}

	// Parse/execute it as a Go template
	out := bytes.Buffer{}
	tpl, err := template.New(normalizedUrl.String()).Parse(string(content))
	if err != nil {
		return
	}

	tpl.Execute(&out, data)

	// Unmarshal the resulting YAML
	err = yaml.Unmarshal(out.Bytes(), &env)
	if err != nil {
		return
	}

	// Process imports if any
	err = processYamlImports(logger, baseLocation, &env, data)
	if err != nil {
		return
	}

	return
}

func processYamlImports(logger *log.Logger, base *url.URL, env *yamlEnvironment, data map[string]interface{}) error {
	if len(env.Imports) > 0 {
		for _, val := range env.Imports {
			importUrl, err := url.Parse(val)
			if err != nil {
				return err
			}
			ref := base.ResolveReference(importUrl)
			logger.Println("Processing import", ref)
			importedDesc, err := parseYamlDescriptor(logger, ref, data)
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
