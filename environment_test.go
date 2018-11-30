package model

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEngineComplete(t *testing.T) {
	env, e := CreateEnvironment(buildUrl("./testdata/yaml/complete.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	assertEnv(t, env)
}

func TestCreateEnginePartials(t *testing.T) {
	env, e := CreateEnvironment(buildUrl("./testdata/yaml/partials/env.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	env2, e := CreateEnvironment(buildUrl("./testdata/yaml/partials/core.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	env.Merge(env2)
	env3, e := CreateEnvironment(buildUrl("./testdata/yaml/partials/providers.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	env.Merge(env3)
	env4, e := CreateEnvironment(buildUrl("./testdata/yaml/partials/orchestrator.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	env.Merge(env4)
	env5, e := CreateEnvironment(buildUrl("./testdata/yaml/partials/stacks.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	env.Merge(env5)
	env6, e := CreateEnvironment(buildUrl("./testdata/yaml/partials/tasks.yaml"), map[string]interface{}{})
	assert.Nil(t, e)
	env.Merge(env6)
	assertEnv(t, env)
}

func assertEnv(t *testing.T, env Environment) {
	assert.Equal(t, "testEnvironment", env.Name)
	assert.Equal(t, "testQualifier", env.Qualifier)
	assert.Equal(t, "This is my awesome Ekara environment.", env.Description)

	// Platform
	assert.NotNil(t, env.Ekara)
	assert.NotNil(t, env.Ekara.Components)
	assert.Equal(t, 5, len(env.Ekara.Components))
	assert.Equal(t, "file://someBase/", env.Ekara.Base.String())
	assert.Equal(t, "file:///someBase/ekara-platform/distribution", env.Ekara.Distribution.Repository.String())

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
	awsComponent, e := providers["aws"].Component.Resolve()
	assert.Nil(t, e)
	assert.True(t, strings.HasSuffix(awsComponent.Repository.String(), "/someBase/ekara-platform/aws-provider"))
	assert.Equal(t, "1.2.3", awsComponent.Version.String())
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
	azureComponent, e := providers["azure"].Component.Resolve()
	assert.Nil(t, e)
	assert.True(t, strings.HasSuffix(azureComponent.Repository.String(), "/someBase/ekara-platform/azure-provider"))
	assert.Equal(t, "1.2.3", azureComponent.Version.String())
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

	//------------------------------------------------------------
	// Node1
	//------------------------------------------------------------
	assert.Equal(t, 10, nodeSets["node1"].Instances)
	ns1Provider, e := nodeSets["node1"].Provider.Resolve()
	assert.Nil(t, e)
	assert.Equal(t, "aws", ns1Provider.Name)

	c = ns1Provider.Parameters
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

	ns1Orchestrator, e := nodeSets["node1"].Orchestrator.Resolve()
	assert.Nil(t, e)
	c = ns1Orchestrator.Parameters
	v, ok = c["orchestrator_node1_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node1_param_key1_value")

	v, ok = c["orchestrator_node1_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node1_param_key2_value")

	c = ns1Orchestrator.Docker
	v, ok = c["docker_node1_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "docker_node1_param_key1_value")

	v, ok = c["docker_node1_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "docker_node1_param_key2_value")

	en = ns1Orchestrator.EnvVars
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

	vol := vs["some/volume/path"]
	assert.Equal(t, vol.Path, "some/volume/path")
	assert.Equal(t, vol.Parameters["param1_name"], "aws_param1_name_value")

	vol = vs["other/volume/path"]
	assert.Equal(t, vol.Path, "other/volume/path")
	assert.Equal(t, vol.Parameters["param2_name"], "aws_param2_name_value")

	//------------------------------------------------------------
	// Node1 Hook
	//------------------------------------------------------------
	no := nodeSets["node1"]
	assert.Equal(t, 1, len(no.Hooks.Provision.Before))
	assert.Equal(t, 1, len(no.Hooks.Provision.After))
	assert.Equal(t, 1, len(no.Hooks.Destroy.Before))
	assert.Equal(t, 1, len(no.Hooks.Destroy.After))

	assert.Equal(t, "task1", no.Hooks.Provision.Before[0].ref)
	assert.Equal(t, "task2", no.Hooks.Provision.After[0].ref)

	assert.Equal(t, "task1", no.Hooks.Destroy.Before[0].ref)
	assert.Equal(t, "task2", no.Hooks.Destroy.After[0].ref)

	//------------------------------------------------------------
	// Node2
	//------------------------------------------------------------
	assert.Equal(t, 20, nodeSets["node2"].Instances)
	ns2Provider, e := nodeSets["node2"].Provider.Resolve()
	assert.Nil(t, e)
	assert.Equal(t, "azure", ns2Provider.Name)

	c = ns2Provider.Parameters
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

	ns2Orchestrator, e := nodeSets["node2"].Orchestrator.Resolve()
	assert.Nil(t, e)
	c = ns2Orchestrator.Parameters
	v, ok = c["orchestrator_node2_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node2_param_key1_value")

	v, ok = c["orchestrator_node2_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node2_param_key2_value")

	c = ns2Orchestrator.Docker
	v, ok = c["docker_node2_param_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "docker_node2_param_key1_value")

	v, ok = c["docker_node2_param_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "docker_node2_param_key2_value")

	en = ns2Orchestrator.EnvVars
	v, ok = en["orchestrator_node2_env_key1"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node2_env_key1_value")

	v, ok = en["orchestrator_node2_env_key2"]
	assert.True(t, ok)
	assert.Equal(t, v, "orchestrator_node2_env_key2_value")

	vs = nodeSets["node2"].Volumes
	assert.NotNil(t, vs)
	assert.Equal(t, 2, len(vs))

	vol = vs["some/volume/path"]
	assert.Equal(t, vol.Path, "some/volume/path")
	assert.Equal(t, vol.Parameters["param1_name"], "azure_param1_name_value")

	vol = vs["other/volume/path"]
	assert.Equal(t, vol.Path, "other/volume/path")
	assert.Equal(t, vol.Parameters["param2_name"], "azure_param2_name_value")

	la = nodeSets["node2"].Labels
	v, ok = la["node2_label1"]
	assert.True(t, ok)
	assert.Equal(t, v, "node2_label1_value")

	v, ok = la["node2_label2"]
	assert.True(t, ok)
	assert.Equal(t, v, "node2_label2_value")

	//------------------------------------------------------------
	// Node2 Hook
	//------------------------------------------------------------
	no = nodeSets["node2"]
	if assert.Equal(t, 1, len(no.Hooks.Provision.Before)) {
		assert.Equal(t, "task1", no.Hooks.Provision.Before[0].ref)
	}
	if assert.Equal(t, 1, len(no.Hooks.Provision.After)) {
		assert.Equal(t, "task2", no.Hooks.Provision.After[0].ref)
	}
	if assert.Equal(t, 1, len(no.Hooks.Destroy.Before)) {
		assert.Equal(t, "task1", no.Hooks.Destroy.Before[0].ref)
	}
	if assert.Equal(t, 1, len(no.Hooks.Destroy.After)) {
		assert.Equal(t, "task2", no.Hooks.Destroy.After[0].ref)
	}

	//------------------------------------------------------------
	// Environment Stacks
	//------------------------------------------------------------
	stacks := env.Stacks
	assert.NotNil(t, stacks)
	assert.Equal(t, 2, len(stacks))

	assert.Contains(t, stacks, "stack1")
	assert.Contains(t, stacks, "stack2")
	assert.NotContains(t, stacks, "dummy")

	stack1 := stacks["stack1"]
	stack2 := stacks["stack2"]

	st1Component, e := stack1.Component.Resolve()
	assert.Nil(t, e)
	assert.True(t, strings.HasSuffix(st1Component.Repository.String(), "/someBase/some-org/stack1"))
	assert.Equal(t, "1.2.3", st1Component.Version.String())

	st2Component, e := stack2.Component.Resolve()
	assert.Nil(t, e)
	assert.True(t, strings.HasSuffix(st2Component.Repository.String(), "/someBase/some-org/stack2"))
	assert.Equal(t, "1.2.3", st2Component.Version.String())

	//------------------------------------------------------------
	// Stack1 Hook
	//------------------------------------------------------------
	if assert.Equal(t, 1, len(stack1.Hooks.Deploy.Before)) {
		assert.Equal(t, "task1", stack1.Hooks.Deploy.Before[0].ref)
	}
	if assert.Equal(t, 1, len(stack1.Hooks.Deploy.After)) {
		assert.Equal(t, "task2", stack1.Hooks.Deploy.After[0].ref)
	}
	if assert.Equal(t, 1, len(stack1.Hooks.Undeploy.Before)) {
		assert.Equal(t, "task1", stack1.Hooks.Undeploy.Before[0].ref)
	}
	if assert.Equal(t, 1, len(stack1.Hooks.Undeploy.After)) {
		assert.Equal(t, "task2", stack1.Hooks.Undeploy.After[0].ref)
	}

	//------------------------------------------------------------
	// Environment Tasks
	//------------------------------------------------------------
	tasks := env.Tasks
	assert.NotNil(t, tasks)
	assert.Equal(t, 3, len(tasks))

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

	//------------------------------------------------------------
	// Environment Tasks Hooks
	//------------------------------------------------------------
	ta := tasks["task3"]
	if assert.Equal(t, 1, len(ta.Hooks.Execute.Before)) {
		assert.Equal(t, "task1", ta.Hooks.Execute.Before[0].ref)
	}
	if assert.Equal(t, 1, len(ta.Hooks.Execute.After)) {
		assert.Equal(t, "task2", ta.Hooks.Execute.After[0].ref)
	}

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

func ExampleEnvironment_Merge() {
	root := Environment{Name: "RootName", Qualifier: "RootQualifier"}
	other := Environment{Name: "OtherName", Qualifier: "OtherQualifier"}
	root.Merge(other)
	fmt.Println(root.QualifiedName())
	// Outpur: RootName_RootQualifier
}
