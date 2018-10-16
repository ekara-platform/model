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
	env, e := Parse(logger, buildUrl("./testdata/yaml/complete.yaml"))
	assert.Nil(t, e)

	assert.Equal(t, "testEnvironment", env.Name)
	assert.Equal(t, "testQualifier", env.Qualifier)
	assert.Equal(t, "This is my awesome Ekara environment.", env.Description)

	// Platform
	assert.NotNil(t, env.Ekara)
	assert.Equal(t, "file://someBase/", env.Ekara.ComponentBase.String())
	assert.Equal(t, "someRegistry.org", env.Ekara.DockerRegistry.String())
	assert.NotNil(t, env.Ekara.Components)
	assert.True(t, strings.HasSuffix(env.Ekara.Component.Resolve().Repository.String(), "someBase/ekara-platform/core"))
	assert.Equal(t, "", env.Ekara.Component.Resolve().Version.String())

	//------------------------------------------------------------
	// Orchestrator
	//------------------------------------------------------------
	orchestrator := env.Orchestrator
	assert.NotNil(t, orchestrator)
	assert.NotNil(t, orchestrator.Parameters)
	c := orchestrator.Parameters
	v, ok := c["swarm_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "swarm_param_key1_value")

	v, ok = c["swarm_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "swarm_param_key2_value")

	c = orchestrator.Docker
	v, ok = c["docker_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "docker_param_key1_value")

	v, ok = c["docker_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "docker_param_key2_value")

	en := orchestrator.EnvVars
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
	assert.True(t, strings.HasSuffix(providers["aws"].Component.Resolve().Repository.String(), "/someBase/ekara-platform/aws-provider"))
	assert.Equal(t, "v1.2.3", providers["aws"].Component.Resolve().Version.String())
	assert.NotNil(t, providers["aws"].Parameters)
	c = providers["aws"].Parameters
	v, ok = c["aws_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "aws_param_key1_value")

	v, ok = c["aws_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "aws_param_key2_value")

	assert.NotNil(t, providers["aws"].EnvVars)
	en = providers["aws"].EnvVars
	v, ok = en["aws_env_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "aws_env_key1_value")

	v, ok = en["aws_env_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "aws_env_key2_value")

	assert.NotNil(t, providers["aws"].Proxy)
	pr := providers["aws"].Proxy
	assert.Equal(t, pr.Http, "aws_http_proxy")
	assert.Equal(t, pr.Https, "aws_https_proxy")
	assert.Equal(t, pr.NoProxy, "aws_no_proxy")

	// Azure Provider
	assert.NotNil(t, providers["azure"])
	assert.Equal(t, "azure", providers["azure"].Name)
	assert.True(t, strings.HasSuffix(providers["azure"].Component.Resolve().Repository.String(), "/someBase/ekara-platform/azure-provider"))
	assert.Equal(t, "v1.2.3", providers["azure"].Component.Resolve().Version.String())
	assert.NotNil(t, providers["azure"].Parameters)

	c = providers["azure"].Parameters
	v, ok = c["azure_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "azure_param_key1_value")

	v, ok = c["azure_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "azure_param_key2_value")

	assert.NotNil(t, providers["azure"].EnvVars)
	en = providers["azure"].EnvVars
	v, ok = en["azure_env_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "azure_env_key1_value")

	v, ok = en["azure_env_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "azure_env_key2_value")

	assert.NotNil(t, providers["azure"].Proxy)
	pr = providers["azure"].Proxy
	assert.Equal(t, pr.Http, "azure_http_proxy")
	assert.Equal(t, pr.Https, "azure_https_proxy")
	assert.Equal(t, pr.NoProxy, "azure_no_proxy")

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

	c = nodeSets["node1"].Provider.Resolve().Parameters
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

	c = nodeSets["node1"].Orchestrator.Resolve().Parameters
	v, ok = c["orchestrator_node1_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node1_param_key1_value")

	v, ok = c["orchestrator_node1_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node1_param_key2_value")

	c = nodeSets["node1"].Orchestrator.Resolve().Docker
	v, ok = c["docker_node1_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "docker_node1_param_key1_value")

	v, ok = c["docker_node1_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "docker_node1_param_key2_value")

	en = nodeSets["node1"].Orchestrator.Resolve().EnvVars
	v, ok = en["orchestrator_node1_env_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node1_env_key1_value")

	v, ok = en["orchestrator_node1_env_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node1_env_key2_value")

	vs := nodeSets["node1"].Volumes
	assert.NotNil(t, vs)
	assert.Equal(t, 2, len(vs))

	la := nodeSets["node1"].Labels
	v, ok = la["node1_label1"]
	assert.True(t, ok)
	assert.Equal(t, v, "node1_label1_value")

	v, ok = la["node1_label2"]
	assert.True(t, ok)
	assert.Equal(t, v, "node1_label2_value")

	vol := vs[0]
	assert.Equal(t, vol.Name, "some/volume/path")
	assert.Equal(t, vol.Parameters["param1_name"], "aws_param1_name_value")

	vol = vs[1]
	assert.Equal(t, vol.Name, "other/volume/path")
	assert.Equal(t, vol.Parameters["param2_name"], "aws_param2_name_value")

	assert.Equal(t, 20, nodeSets["node2"].Instances)
	assert.Equal(t, "azure", nodeSets["node2"].Provider.provider.Name)

	c = nodeSets["node2"].Provider.Resolve().Parameters
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

	c = nodeSets["node2"].Orchestrator.Resolve().Parameters
	v, ok = c["orchestrator_node2_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node2_param_key1_value")

	v, ok = c["orchestrator_node2_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node2_param_key2_value")

	c = nodeSets["node2"].Orchestrator.Resolve().Docker
	v, ok = c["docker_node2_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "docker_node2_param_key1_value")

	v, ok = c["docker_node2_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "docker_node2_param_key2_value")

	en = nodeSets["node2"].Orchestrator.Resolve().EnvVars
	v, ok = en["orchestrator_node2_env_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node2_env_key1_value")

	v, ok = en["orchestrator_node2_env_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node2_env_key2_value")

	vs = nodeSets["node2"].Volumes
	assert.NotNil(t, vs)
	assert.Equal(t, 2, len(vs))

	vol = vs[0]
	assert.Equal(t, vol.Name, "some/volume/path")
	assert.Equal(t, vol.Parameters["param1_name"], "azure_param1_name_value")

	vol = vs[1]
	assert.Equal(t, vol.Name, "other/volume/path")
	assert.Equal(t, vol.Parameters["param2_name"], "azure_param2_name_value")

	la = nodeSets["node2"].Labels
	v, ok = la["node2_label1"]
	assert.True(t, ok)
	assert.Equal(t, v, "node2_label1_value")

	v, ok = la["node2_label2"]
	assert.True(t, ok)
	assert.Equal(t, v, "node2_label2_value")

	//------------------------------------------------------------
	// Environment Stacks
	//------------------------------------------------------------
	stacks := env.Stacks
	assert.NotNil(t, stacks)
	assert.Equal(t, 2, len(stacks))

	assert.Contains(t, stacks, "stack1")
	assert.Contains(t, stacks, "stack2")
	assert.NotContains(t, stacks, "dummy")

	assert.True(t, strings.HasSuffix(stacks["stack1"].Component.Resolve().Repository.String(), "/someBase/some-org/stack1"))
	assert.Equal(t, "v1.2.3", stacks["stack1"].Component.Resolve().Version.String())

	assert.True(t, strings.HasSuffix(stacks["stack2"].Component.Resolve().Repository.String(), "/someBase/some-org/stack2"))
	assert.Equal(t, "v1.2.3", stacks["stack2"].Component.Resolve().Version.String())

	//------------------------------------------------------------
	// Environment Tasks
	//------------------------------------------------------------
	tasks := env.Tasks
	assert.NotNil(t, tasks)
	assert.Equal(t, 2, len(tasks))

	assert.Contains(t, tasks, "task1")

	pa := tasks["task1"].Parameters
	v, ok = pa["tasks_task1_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task1_param_key1_value")

	v, ok = pa["tasks_task1_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task1_param_key2_value")

	en = tasks["task1"].EnvVars
	v, ok = en["tasks_task1_env_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task1_env_key1_value")

	v, ok = en["tasks_task1_env_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task1_env_key2_value")

	assert.Contains(t, tasks, "task2")

	pa = tasks["task2"].Parameters
	v, ok = pa["tasks_task2_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task2_param_key1_value")

	v, ok = pa["tasks_task2_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task2_param_key2_value")

	en = tasks["task2"].EnvVars
	v, ok = en["tasks_task2_env_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task2_env_key1_value")

	v, ok = en["tasks_task2_env_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "tasks_task2_env_key2_value")

	assert.NotContains(t, tasks, "dummy")

	assert.Equal(t, "task1_playbook", tasks["task1"].Playbook)
	assert.Equal(t, "task1_cron", tasks["task1"].Cron)

	assert.Equal(t, "task2_playbook", tasks["task2"].Playbook)
	assert.Equal(t, "task2_cron", tasks["task2"].Cron)
}

func buildUrl(loc string) *url.URL {
	u, e := url.Parse(loc)
	if e != nil {
		panic(e)
	}
	return u
}

func TestQualifiedName(t *testing.T) {
	env := Environment{
		Name:      "MyName",
		Qualifier: "MyQualifier",
	}
	qn := env.QualifiedName()
	assert.NotNil(t, qn)
	assert.Equal(t, "MyName_MyQualifier", qn.String())
}

func TestUnqualifiedName(t *testing.T) {
	env := Environment{
		Name: "MyName",
	}
	qn := env.QualifiedName()
	assert.NotNil(t, qn)
	assert.Equal(t, "MyName", qn.String())
}
