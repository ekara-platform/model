package model

import (
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"log"
)

type yamlLabels struct {
	Labels []string
}

type yamlParams struct {
	Params map[string]string
}

type yamlRef struct {
	yamlParams `yaml:",inline"`
	Name       string
}

type yamlHook struct {
	Before []yamlRef
	After  []yamlRef
}

type yamlEnvironment struct {
	yamlLabels `yaml:",inline"`

	// Global attributes
	Name        string
	Description string
	Version     string

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

	// Imports
	Imports []string

	// Components
	Components map[string]string

	// Providers
	Providers map[string]struct {
		yamlParams `yaml:",inline"`

		Repository string
		Version    string
	}

	// Node sets
	Nodes map[string]struct {
		yamlLabels `yaml:",inline"`

		Provider  yamlRef
		Instances int
		Hooks     struct {
			Provision yamlHook
			Destroy   yamlHook
		}
	}

	// Software stacks
	Stacks map[string]struct {
		yamlLabels `yaml:",inline"`

		Repository string
		Version    string
		DeployOn   []string `yaml:"deployOn"`
		Hooks      struct {
			Deploy   yamlHook
			Undeploy yamlHook
		}
	}

	// Custom tasks
	Tasks map[string]struct {
		yamlLabels `yaml:",inline"`
		yamlParams `yaml:",inline"`

		Playbook string
		Cron     string
		RunOn    []string `yaml:"runOn"`
		Hooks    struct {
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

func parseYamlDescriptor(logger *log.Logger, location string) (env yamlEnvironment, err error) {
	baseLocation, content, err := readUrlOrFile(logger, location)
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

func processYamlImports(logger *log.Logger, baseLocation string, env *yamlEnvironment) error {
	if len(env.Imports) > 0 {
		logger.Println("Processing imports", env.Imports)
		for _, val := range env.Imports {
			logger.Println("Processing import", baseLocation+val)
			importedDesc, err := parseYamlDescriptor(logger, baseLocation+val)
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
