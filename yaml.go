package model

import (
	"log"

	"net/url"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

type yamlLabels struct {
	Labels []string
}

type yamlParams struct {
	Params attributes
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

	// Orchestrator
	Orchestrator struct {
		yamlParams `yaml:",inline"`

		Name       string
		Repository string
		Version    string
	}

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

		// Orchestrator
		Orchestrator struct {
			yamlParams `yaml:",inline"`
		}

		Hooks struct {
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

func parseYamlDescriptor(logger *log.Logger, u *url.URL) (env yamlEnvironment, err error) {
	normalizedUrl, err := NormalizeUrl(u)
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
