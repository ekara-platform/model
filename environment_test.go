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

	// Platform
	assert.NotNil(t, env.Lagoon)
	assert.Equal(t, "file://someBase/", env.Lagoon.ComponentBase.String())
	assert.Equal(t, "someRegistry.org", env.Lagoon.DockerRegistry.String())
	assert.Equal(t, "http://user:pwd@someproxy.org:8080", env.Lagoon.Proxy.Http.String())
	assert.Equal(t, "https://user:pwd@someproxy.org:8080", env.Lagoon.Proxy.Https.String())
	assert.Equal(t, "*.dummy.org", env.Lagoon.Proxy.NoProxy)
	assert.NotNil(t, env.Lagoon.ComponentVersions)
	assert.True(t, strings.HasSuffix(env.Lagoon.Component.Repository.String(), "someBase/lagoon-platform/core"))
	assert.Equal(t, "stable", env.Lagoon.Component.Version.String())

	//------------------------------------------------------------
	// Orchestrator
	//------------------------------------------------------------
	orchestrator := env.Orchestrator
	assert.NotNil(t, orchestrator)
	assert.NotNil(t, orchestrator.Parameters)
	c := orchestrator.Parameters.copy()
	v, ok := c["swarm_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "swarm_param_key1_value")

	v, ok = c["swarm_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "swarm_param_key2_value")

	en := orchestrator.Envvars.copy()
	v, ok = en["swarm_env_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "swarm_env_key1_value")

	v, ok = en["swarm_env_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "swarm_env_key2_value")

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

	assert.NotNil(t, providers["aws"].Envvars)
	en = providers["aws"].Envvars.copy()
	v, ok = en["aws_env_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "aws_env_key1_value")

	v, ok = en["aws_env_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "aws_env_key2_value")

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

	assert.NotNil(t, providers["azure"].Envvars)
	en = providers["azure"].Envvars.copy()
	v, ok = en["azure_env_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "azure_env_key1_value")

	v, ok = en["azure_env_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "azure_env_key2_value")

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

	en = nodeSets["node1"].Orchestrator.Envvars.copy()
	v, ok = en["orchestrator_node1_env_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node1_env_key1_value")

	v, ok = en["orchestrator_node1_env_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node1_env_key2_value")

	vs := nodeSets["node1"].Provider.Volumes()
	assert.NotNil(t, vs)
	assert.Equal(t, 2, len(vs))

	vol := vs[0]
	assert.Equal(t, vol.Name, "aws_name1")
	assert.Equal(t, vol.Parameters["param1_name"], "aws_param1_name_value")

	vol = vs[1]
	assert.Equal(t, vol.Name, "aws_name2")
	assert.Equal(t, vol.Parameters["param2_name"], "aws_param2_name_value")

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

	en = nodeSets["node2"].Orchestrator.Envvars.copy()
	v, ok = en["orchestrator_node2_env_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node2_env_key1_value")

	v, ok = en["orchestrator_node2_env_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node2_env_key2_value")

	vs = nodeSets["node2"].Provider.Volumes()
	assert.NotNil(t, vs)
	assert.Equal(t, 2, len(vs))

	vol = vs[0]
	assert.Equal(t, vol.Name, "azure_name1")
	assert.Equal(t, vol.Parameters["param1_name"], "azure_param1_name_value")

	vol = vs[1]
	assert.Equal(t, vol.Name, "azure_name2")
	assert.Equal(t, vol.Parameters["param2_name"], "azure_param2_name_value")

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

	pa := tasks["task1"].Parameters.copy()
	v, ok = pa["tasks_task1_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task1_param_key1_value")

	v, ok = pa["tasks_task1_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task1_param_key2_value")

	en = tasks["task1"].Envvars.copy()
	v, ok = en["tasks_task1_env_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task1_env_key1_value")

	v, ok = en["tasks_task1_env_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task1_env_key2_value")

	assert.Contains(t, tasks, "task2")

	pa = tasks["task2"].Parameters.copy()
	v, ok = pa["tasks_task2_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task2_param_key1_value")

	v, ok = pa["tasks_task2_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task2_param_key2_value")

	en = tasks["task2"].Envvars.copy()
	v, ok = en["tasks_task2_env_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task2_env_key1_value")

	v, ok = en["tasks_task2_env_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task2_env_key2_value")

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
