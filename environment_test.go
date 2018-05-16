package model

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEngineComplete(t *testing.T) {
	logger := log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime)
	env, e := Parse(logger, "./testdata/yaml/complete_descriptor/lagoon.yaml")
	assert.Nil(t, e)

	assert.Equal(t, "name_value", env.Name)
	assert.Equal(t, "description_value", env.Description)

	// Environment Version
	assert.Equal(t, 1, env.Version.Major)
	assert.Equal(t, 2, env.Version.Minor)
	assert.Equal(t, 3, env.Version.Micro)
	assert.Equal(t, "1.2.3", env.Version.Full)

	// Environment Labels
	labels := env.Labels
	assert.NotNil(t, labels)
	assert.Equal(t, 3, len(labels.AsStrings()))
	assert.Equal(t, true, labels.MatchesLabels("root_label1", "root_label2", "root_label3"))

	// Settings
	assert.NotNil(t, env.Settings)
	assert.Equal(t, "http://user:pwd@someproxy.org:8080", env.Settings.Proxy.Http.String())
	assert.Equal(t, "https://user:pwd@someproxy.org:8080", env.Settings.Proxy.Https.String())
	assert.Equal(t, "*.dummy.org", env.Settings.Proxy.NoProxy)

	//------------------------------------------------------------
	// Components
	//------------------------------------------------------------
	components := env.Components
	assert.NotNil(t, components)

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
	assert.Equal(t, "https://github.com/lagoon-platform/aws-provider.git", providers["aws"].Repository.String())
	assert.Equal(t, "1.2.3", providers["aws"].Version.Full)
	assert.NotNil(t, providers["aws"].Parameters)
	assert.Equal(t, map[string]string{"aws_param_key1": "aws_param_key1_value", "aws_param_key2": "aws_param_key2_value"}, providers["aws"].Parameters.AsMap())

	// Azure Provider
	assert.NotNil(t, providers["azure"])
	assert.Equal(t, "azure", providers["azure"].Name)
	assert.Equal(t, "https://github.com/lagoon-platform/azure-provider.git", providers["azure"].Repository.String())
	assert.Equal(t, "1.2.3", providers["azure"].Version.Full)
	assert.NotNil(t, providers["azure"].Parameters)
	assert.Equal(t, map[string]string{"azure_param_key1": "azure_param_key1_value", "azure_param_key2": "azure_param_key2_value"}, providers["azure"].Parameters.AsMap())

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
	assert.Equal(t, "aws", nodeSets["node1"].Provider.Resolve().Name)
	assert.Equal(t, []string{"node1_label1", "node1_label2", "node1_label3"}, nodeSets["node1"].Labels.AsStrings())
	assert.Equal(t, map[string]string{
		"aws_param_key1":            "aws_param_key1_value",
		"aws_param_key2":            "aws_param_key2_value",
		"provider_node1_param_key1": "provider_node1_param_key1_value",
		"provider_node1_param_key2": "provider_node1_param_key2_value"},
		nodeSets["node1"].Provider.Resolve().Parameters.AsMap())

	assert.Equal(t, 20, nodeSets["node2"].Instances)
	assert.Equal(t, []string{"node2_label1", "node2_label2", "node2_label3"}, nodeSets["node2"].Labels.AsStrings())
	assert.Equal(t, "azure", nodeSets["node2"].Provider.Resolve().Name)
	assert.Equal(t, map[string]string{
		"azure_param_key1":          "azure_param_key1_value",
		"azure_param_key2":          "azure_param_key2_value",
		"provider_node2_param_key1": "provider_node2_param_key1_value",
		"provider_node2_param_key2": "provider_node2_param_key2_value"},
		nodeSets["node2"].Provider.Resolve().Parameters.AsMap())

	//------------------------------------------------------------
	// Environment Stacks
	//------------------------------------------------------------
	stacks := env.Stacks
	assert.NotNil(t, stacks)
	assert.Equal(t, 2, len(stacks))

	assert.Contains(t, stacks, "stack1")
	assert.Contains(t, stacks, "stack2")
	assert.NotContains(t, stacks, "dummy")

	assert.Equal(t, "https://github.com/lagoon-platform/stack1_repository.git", stacks["stack1"].Repository.String())
	assert.Equal(t, "1.2.3", stacks["stack1"].Version.Full)
	assert.Equal(t, []string{"stack1_label1", "stack1_label2", "stack1_label3"}, stacks["stack1"].Labels.AsStrings())

	assert.Equal(t, "https://github.com/lagoon-platform/stack2_repository.git", stacks["stack2"].Repository.String())
	assert.Equal(t, "1.2.3", stacks["stack2"].Version.Full)
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
