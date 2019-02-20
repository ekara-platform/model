package model

import (
	"net/url"

	"bytes"
	"text/template"

	"gopkg.in/yaml.v2"
)

type (
	// yaml tag for the proxy details
	yamlProxy struct {
		Http    string `yaml:"http_proxy"`
		Https   string `yaml:"https_proxy"`
		NoProxy string `yaml:"no_proxy"`
	}

	// yaml tag for parameters
	yamlParams struct {
		Params map[string]interface{} `yaml:",omitempty"`
	}

	// yaml tag for Docker parameters
	yamlDockerParams struct {
		Docker map[string]interface{} `yaml:",omitempty"`
	}

	// yaml tag for authentication parameters
	yamlAuth struct {
		Auth map[string]interface{} `yaml:",omitempty"`
	}

	// yaml tag for environment variables
	yamlEnv struct {
		Env map[string]string `yaml:",omitempty"`
	}

	// yaml tag for labels on nodesets
	yamlLabel struct {
		Labels map[string]string `yaml:",omitempty"`
	}

	// yaml tag for component
	yamlComponent struct {
		// The source repository where the component lives
		Repository string
		// The ref (branch or tag) of the component to use
		Ref string
		// The authentication parameters
		yamlAuth `yaml:",inline"`
	}

	// yaml tag for component with imports
	yamlComponentWithImports struct {
		yamlComponent `yaml:",inline"`
		// Local imports for the component
		Imports []string `yaml:",omitempty"`
	}

	// yaml tag for a volume and its parameters
	yamlVolume struct {
		// The mounting path of the created volume
		Path string
		// The parameters required to create the volume (typically provider dependent)
		yamlParams `yaml:",inline"`
	}

	// yaml tag for a shared volume content
	yamlVolumeContent struct {
		// The component holding the content to copy into the volume
		Component string
		// The path of the content to copy
		Path string
	}

	// yaml reference to provider
	yamlProviderRef struct {
		Name string
		// The overriding provider parameters
		yamlParams `yaml:",inline"`
		// The overriding provider environment variables
		yamlEnv `yaml:",inline"`
		// The overriding provider proxy
		Proxy yamlProxy
	}

	// yaml reference to orchestrator
	yamlOrchestratorRef struct {
		// The overriding orchestrator parameters
		yamlParams `yaml:",inline"`
		// The overriding docker parameters
		yamlDockerParams `yaml:",inline"`
		// The overriding orchestrator environment variables
		yamlEnv `yaml:",inline"`
	}

	// yaml reference to task
	yamlTaskRef struct {
		// The referenced task
		Task string
		// The overriding parameters
		yamlParams `yaml:",inline"`
		// The overriding environment variables
		yamlEnv `yaml:",inline"`
	}

	//yaml tag for hooks
	yamlHook struct {
		// Hooks to be executed before the corresponding process step
		Before []yamlTaskRef `yaml:",omitempty"`
		// Hooks to be executed after the corresponding process step
		After []yamlTaskRef `yaml:",omitempty"`
	}

	// Definition of the Ekara environment
	yamlEnvironment struct {
		// The name of the environment
		Name string
		// The qualifier of the environment
		Qualifier string `yaml:",omitempty"`

		// The description of the environment
		Description string `yaml:",omitempty"`

		// The Ekara platform used to interact with the environment
		Ekara struct {
			Base         string `yaml:",omitempty"`
			Distribution yamlComponent
			Components   map[string]yamlComponentWithImports
		}

		// Global imports
		Imports []string `yaml:",omitempty"`

		// Tasks which can be run on the created environment
		Tasks map[string]struct {
			// Name of the task component
			Component string
			// The task parameters
			yamlParams `yaml:",inline"`
			// The task environment variables
			yamlEnv `yaml:",inline"`
			// The name of the playbook to launch the task
			Playbook string `yaml:",omitempty"`
			// The CRON to run cyclically the task
			Cron string `yaml:",omitempty"`
			// The Hooks to be executed in addition the the main task playbook
			Hooks struct {
				Execute yamlHook `yaml:",omitempty"`
			} `yaml:",omitempty"`
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
				Provision yamlHook `yaml:",omitempty"`
				Destroy   yamlHook `yaml:",omitempty"`
			} `yaml:",omitempty"`

			// The labels associated with the nodeset
			yamlLabel `yaml:",inline"`
		}

		// Software stacks to be installed on the environment
		Stacks map[string]struct {
			// Name of the stack component
			Component string
			// The Hooks to be executed while deploying and undeploying the stack
			Hooks struct {
				Deploy   yamlHook `yaml:",omitempty"`
				Undeploy yamlHook `yaml:",omitempty"`
			} `yaml:",omitempty"`

			// The parameters
			yamlParams `yaml:",inline"`
			// The environment variables
			yamlEnv `yaml:",inline"`
		}

		// Global hooks
		Hooks struct {
			Init      yamlHook `yaml:",omitempty"`
			Provision yamlHook `yaml:",omitempty"`
			Deploy    yamlHook `yaml:",omitempty"`
			Undeploy  yamlHook `yaml:",omitempty"`
			Destroy   yamlHook `yaml:",omitempty"`
		} `yaml:",omitempty"`

		// Global volumes
		Volumes map[string]struct {
			Content []yamlVolumeContent `yaml:",omitempty"`
		} `yaml:",omitempty"`
	}
)

// RawContent returns the serialized content of the environement as YAML
func (r *yamlEnvironment) RawContent() ([]byte, error) {
	return yaml.Marshal(r)
}

func parseYamlDescriptor(u *url.URL, data map[string]interface{}) (env yamlEnvironment, err error) {
	var normalizedUrl *url.URL
	normalizedUrl, err = NormalizeUrl(u)
	if err != nil {
		return
	}

	// Read descriptor content
	content, err := ReadUrl(normalizedUrl)
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

	return
}
