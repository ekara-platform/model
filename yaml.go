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
		Params map[string]interface{}
	}

	// yaml tag for Docker parameters
	yamlDockerParams struct {
		Docker map[string]interface{}
	}

	// yaml tag for authentication parameters
	yamlAuth struct {
		Auth map[string]interface{}
	}

	// yaml tag for environment variables
	yamlEnv struct {
		Env map[string]string
	}

	// yaml tag for labels on nodesets
	yamlLabel struct {
		Labels map[string]string
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
		Imports []string
	}

	// yaml tag for a volume and its parameters
	yamlVolume struct {
		// The mounting path of the created volume
		Path string
		// The parameters required to create the volume (typically provider dependent)
		yamlParams `yaml:",inline"`
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
		Before []yamlTaskRef
		// Hooks to be executed after the corresponding process step
		After []yamlTaskRef
	}

	// Definition of the Ekara environment
	yamlEnvironment struct {
		// The name of the environment
		Name string
		// The qualifier of the environment
		Qualifier string

		// The description of the environment
		Description string

		// The Ekara platform used to interact with the environment
		Ekara struct {
			Base         string
			Distribution yamlComponent
			Components   map[string]yamlComponentWithImports
		}

		// Global imports
		Imports []string

		// Tasks which can be run on the created environment
		Tasks map[string]struct {
			// Name of the task component
			Component string
			// The task parameters
			yamlParams `yaml:",inline"`
			// The task environment variables
			yamlEnv `yaml:",inline"`
			// The name of the playbook to launch the task
			Playbook string
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
			// The labels associated with the nodeset
			yamlLabel `yaml:",inline"`
		}

		// Software stacks to be installed on the environment
		Stacks map[string]struct {
			// Name of the stack component
			Component string
			// The Hooks to be executed while deploying and undeploying the stack
			Hooks struct {
				Deploy   yamlHook
				Undeploy yamlHook
			}
			// The parameters
			yamlParams `yaml:",inline"`
			// The environment variables
			yamlEnv `yaml:",inline"`
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
