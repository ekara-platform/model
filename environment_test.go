package model

import (
	"log"
	"os"
	"testing"

	"net/url"
	"strings"

	"github.com/stretchr/testify/assert"
)

func TestCreateEngineComplete(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := Parse(logger, buildUrl("./testdata/yaml/complete_descriptor/lagoon.yaml"))
	assert.Nil(t, e)

	assert.Equal(t, "name_value", env.Name)
	assert.Equal(t, "description_value", env.Description)

	// Environment Version
	assert.Equal(t, 1, env.Version.Major)
	assert.Equal(t, 2, env.Version.Minor)
	assert.Equal(t, 3, env.Version.Micro)
	assert.Equal(t, "v1.2.3", env.Version.String())

	// Environment Labels
	labels := env.Labels
	assert.NotNil(t, labels)
	assert.Equal(t, 3, len(labels.AsStrings()))
	assert.Equal(t, true, labels.MatchesLabels("root_label1", "root_label2", "root_label3"))

	// Settings
	assert.NotNil(t, env.Settings)
	assert.Equal(t, "file://someBase/", env.Settings.ComponentBase.String())
	assert.Equal(t, "someRegistry.org", env.Settings.DockerRegistry.String())
	assert.Equal(t, "http://user:pwd@someproxy.org:8080", env.Settings.Proxy.Http.String())
	assert.Equal(t, "https://user:pwd@someproxy.org:8080", env.Settings.Proxy.Https.String())
	assert.Equal(t, "*.dummy.org", env.Settings.Proxy.NoProxy)

	//------------------------------------------------------------
	// Components
	//------------------------------------------------------------
	components := env.Components
	assert.NotNil(t, components)

	//------------------------------------------------------------
	// Orchestrator
	//------------------------------------------------------------
	orchestrator := env.Orchestrator
	assert.NotNil(t, orchestrator)
	assert.Equal(t, "swarm", orchestrator.Name)

	assert.NotNil(t, orchestrator.Parameters)
	c := orchestrator.Parameters.copy()
	v, ok := c["swarm_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "swarm_param_key1_value")

	v, ok = c["swarm_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "swarm_param_key2_value")

	//------------------------------------------------------------
	// Environment Providers
	//------------------------------------------------------------
	providers := env.Providers
	assert.NotNil(t, providers)
	assert.Equal(t, 2, len(providers))

	assert.Contains(t, providers, "aws")
	assert.Contains(t, providers, "azure")
	assert.NotContains(t, providers, "dummy")

	// AWS Provider
	assert.NotNil(t, providers["aws"])
	assert.Equal(t, "aws", providers["aws"].Name)
	assert.True(t, strings.HasSuffix(providers["aws"].Repository.String(), "/someBase/lagoon-platform/aws-provider"))
	assert.Equal(t, "v1.2.3", providers["aws"].Version.String())
	assert.NotNil(t, providers["aws"].Parameters)
	c = providers["aws"].Parameters.copy()
	v, ok = c["aws_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "aws_param_key1_value")

	v, ok = c["aws_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "aws_param_key2_value")

	// Azure Provider
	assert.NotNil(t, providers["azure"])
	assert.Equal(t, "azure", providers["azure"].Name)
	assert.True(t, strings.HasSuffix(providers["azure"].Repository.String(), "/someBase/lagoon-platform/azure-provider"))
	assert.Equal(t, "v1.2.3", providers["azure"].Version.String())
	assert.NotNil(t, providers["azure"].Parameters)

	c = providers["azure"].Parameters.copy()
	v, ok = c["azure_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "azure_param_key1_value")

	v, ok = c["azure_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "azure_param_key2_value")

	//------------------------------------------------------------
	// Environment Nodes
	//------------------------------------------------------------
	nodeSets := env.NodeSets
	assert.NotNil(t, nodeSets)
	assert.Equal(t, 2, len(nodeSets))

	assert.Contains(t, nodeSets, "node1")
	assert.Contains(t, nodeSets, "node2")
	assert.NotContains(t, nodeSets, "dummy")

	assert.Equal(t, 10, nodeSets["node1"].Instances)
	assert.Equal(t, "aws", nodeSets["node1"].Provider.provider.Name)
	assert.Equal(t, []string{"node1_label1", "node1_label2", "node1_label3"}, nodeSets["node1"].Labels.AsStrings())

	c = nodeSets["node1"].Provider.Parameters.copy()
	v, ok = c["aws_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "aws_param_key1_value")

	v, ok = c["aws_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "aws_param_key2_value")

	v, ok = c["provider_node1_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "provider_node1_param_key1_value")

	v, ok = c["provider_node1_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "provider_node1_param_key2_value")

	c = nodeSets["node1"].Orchestrator.Parameters.copy()
	v, ok = c["orchestrator_node1_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node1_param_key1_value")

	v, ok = c["orchestrator_node1_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node1_param_key2_value")

	assert.Equal(t, 20, nodeSets["node2"].Instances)
	assert.Equal(t, []string{"node2_label1", "node2_label2", "node2_label3"}, nodeSets["node2"].Labels.AsStrings())
	assert.Equal(t, "azure", nodeSets["node2"].Provider.provider.Name)

	c = nodeSets["node2"].Provider.Parameters.copy()
	v, ok = c["azure_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "azure_param_key1_value")

	v, ok = c["azure_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "azure_param_key2_value")

	v, ok = c["provider_node2_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "provider_node2_param_key1_value")

	v, ok = c["provider_node2_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "provider_node2_param_key2_value")

	c = nodeSets["node2"].Orchestrator.Parameters.copy()
	v, ok = c["orchestrator_node2_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node2_param_key1_value")

	v, ok = c["orchestrator_node2_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node2_param_key2_value")

	//------------------------------------------------------------
	// Environment Stacks
	//------------------------------------------------------------
	stacks := env.Stacks
	assert.NotNil(t, stacks)
	assert.Equal(t, 2, len(stacks))

	assert.Contains(t, stacks, "stack1")
	assert.Contains(t, stacks, "stack2")
	assert.NotContains(t, stacks, "dummy")

	assert.True(t, strings.HasSuffix(stacks["stack1"].Repository.String(), "/someBase/lagoon-platform/stack1_repository"))
	assert.Equal(t, "v1.2.3", stacks["stack1"].Version.String())
	assert.Equal(t, []string{"stack1_label1", "stack1_label2", "stack1_label3"}, stacks["stack1"].Labels.AsStrings())

	assert.True(t, strings.HasSuffix(stacks["stack2"].Repository.String(), "/someBase/lagoon-platform/stack2_repository"))
	assert.Equal(t, "v1.2.3", stacks["stack2"].Version.String())
	assert.Equal(t, []string{"stack2_label1", "stack2_label2", "stack2_label3"}, stacks["stack2"].Labels.AsStrings())

	//------------------------------------------------------------
	// Environment Tasks
	//------------------------------------------------------------
	tasks := env.Tasks
	assert.NotNil(t, tasks)
	assert.Equal(t, 2, len(tasks))

	assert.Contains(t, tasks, "task1")
	assert.Contains(t, tasks, "task2")
	assert.NotContains(t, tasks, "dummy")

	assert.Equal(t, "task1_playbook", tasks["task1"].Playbook)
	assert.Equal(t, "task1_cron", tasks["task1"].Cron)
	assert.Equal(t, []string{"task1_label1", "task1_label2", "task1_label3"}, tasks["task1"].Labels.AsStrings())

	assert.Equal(t, "task2_playbook", tasks["task2"].Playbook)
	assert.Equal(t, "task2_cron", tasks["task2"].Cron)
	assert.Equal(t, []string{"task2_label1", "task2_label2", "task2_label3"}, tasks["task2"].Labels.AsStrings())
}

func buildUrl(loc string) *url.URL {
	u, e := url.Parse(loc)
	if e != nil {
		panic(e)
	}
	return u
}
