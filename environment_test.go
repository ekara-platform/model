package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEngineComplete(t *testing.T) {

	yamlEnv, e := ParseYamlDescriptor(buildURL(t, "./testdata/yaml/complete.yaml"), &TemplateContext{})
	assert.Nil(t, e)
	p, e := createPlatform(yamlEnv.Ekara)
	assert.Nil(t, e)
	env, e := CreateEnvironment("", yamlEnv, MainComponentId)
	assert.Nil(t, e)
	env.ekara = &p
	assertEnv(t, env)
}

func assertEnv(t *testing.T, env *Environment) {
	assert.Equal(t, "testEnvironment", env.Name)
	assert.Equal(t, "testQualifier", env.Qualifier)
	assert.Equal(t, "This is my awesome Ekara environment.", env.Description)

	// Platform
	assert.NotNil(t, env.ekara)
	assert.NotNil(t, env.ekara.Components)
	assert.Equal(t, 8, len(env.ekara.Components))
	assert.Equal(t, SchemeFile, env.ekara.Base.Url.UpperScheme())
	assert.Equal(t, "file://someBase/", env.ekara.Base.Url.String())
	assert.Equal(t, "file:///someBase/ekara-platform/distribution/", env.ekara.Parent.Repository.Url.String())

	// Variables
	assert.NotNil(t, env.Vars)
	if assert.Equal(t, 2, len(env.Vars)) {
		va, ok := env.Vars["global_var_key1"]
		assert.True(t, ok)
		assert.Equal(t, va, "global_var_val1")

		va, ok = env.Vars["global_var_key2"]
		assert.True(t, ok)
		assert.Equal(t, va, "global_var_val2")
	}

	// Templates
	assert.NotNil(t, env.Templates)
	templates := env.Templates
	tC := templates.Content
	if assert.Equal(t, len(tC), 2) {
		assert.Contains(t, tC, "environment/*/*.yaml")
		assert.Contains(t, tC, "environment/*.yml")
	}

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
	awsComponent, err := providers["aws"].cRef.resolve()
	assert.Nil(t, err)

	assert.True(t, strings.HasSuffix(awsComponent.Repository.Url.String(), "/someBase/ekara-platform/aws-provider/"))
	assert.Equal(t, "1.2.3", awsComponent.Repository.Ref)
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
	azureComponent, err := providers["azure"].cRef.resolve()
	assert.Nil(t, err)
	assert.True(t, strings.HasSuffix(azureComponent.Repository.Url.String(), "/someBase/ekara-platform/azure-provider/"))
	assert.Equal(t, "1.2.3", azureComponent.Repository.Ref)
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

	la := nodeSets["node1"].Labels
	v, ok = la["node1_label1"]
	assert.True(t, ok)
	assert.Equal(t, v, "node1_label1_value")

	v, ok = la["node1_label2"]
	assert.True(t, ok)
	assert.Equal(t, v, "node1_label2_value")

	//------------------------------------------------------------
	// Node1 Hook
	//------------------------------------------------------------
	no := nodeSets["node1"]
	assert.Equal(t, 1, len(no.Hooks.Provision.Before))
	assert.Equal(t, 1, len(no.Hooks.Provision.After))

	assert.Equal(t, "task1", no.Hooks.Provision.Before[0].ref)
	assert.Equal(t, "task2", no.Hooks.Provision.After[0].ref)

	//------------------------------------------------------------
	// Node1 Hook Env and Param
	//------------------------------------------------------------
	r, err := no.Hooks.Provision.After[0].Resolve()
	assert.Nil(t, err)
	p := r.Parameters

	if assert.Equal(t, 3, len(p)) {
		assert.Equal(t, "tasks_task2_param_key1_value_overwritten", p["tasks_task2_param_key1"])
		assert.Equal(t, "tasks_task2_param_key2_value", p["tasks_task2_param_key2"])
		assert.Equal(t, "tasks_task2_param_key3_value", p["tasks_task2_param_key3"])

	}
	envvars := r.EnvVars
	if assert.Equal(t, 3, len(envvars)) {
		assert.Equal(t, "tasks_task2_env_key1_value_overwritten", envvars["tasks_task2_env_key1"])
		assert.Equal(t, "tasks_task2_env_key2_value", envvars["tasks_task2_env_key2"])
		assert.Equal(t, "tasks_task2_env_key3_value", envvars["tasks_task2_env_key3"])
	}

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

	//------------------------------------------------------------
	//Stack denpendency
	//------------------------------------------------------------
	b, sd := stack1.Dependency()
	assert.False(t, b)
	assert.Equal(t, len(sd), 0)

	b, sd = stack2.Dependency()
	assert.True(t, b)
	assert.Equal(t, len(sd), 1)
	assert.Equal(t, stack1.Name, sd[0].Name)

	st1Component, err := stack1.cRef.resolve()
	assert.Nil(t, err)
	assert.True(t, strings.HasSuffix(st1Component.Repository.Url.String(), "/someBase/some-org/stack1/"))
	assert.Equal(t, "1.2.3", st1Component.Repository.Ref)

	st2Component, err := stack2.cRef.resolve()
	assert.Nil(t, err)
	assert.True(t, strings.HasSuffix(st2Component.Repository.Url.String(), "/someBase/some-org/stack2/"))
	assert.Equal(t, "1.2.3", st2Component.Repository.Ref)

	//------------------------------------------------------------
	//Stack copies
	//------------------------------------------------------------
	copies := stack2.Copies
	if assert.Equal(t, len(copies.Content), 2) {
		if assert.Contains(t, copies.Content, "cp1") {
			v, ok := copies.Content["cp1"]
			assert.Equal(t, v.Path, "some/target1/volume/path")
			assert.True(t, v.Once)
			assert.True(t, ok)
			assert.Contains(t, v.Sources.Content, "*target1_to_be_copied.yaml")
			assert.Contains(t, v.Sources.Content, "*target1_to_be_copied.yml")
			lab, ok := v.Labels["label1"]
			assert.True(t, ok)
			assert.Equal(t, lab, "t1_val1")
			lab, ok = v.Labels["label2"]
			assert.True(t, ok)
			assert.Equal(t, lab, "t1_val2")
		}
		if assert.Contains(t, copies.Content, "cp2") {
			v, ok := copies.Content["cp2"]
			assert.Equal(t, v.Path, "some/target2/volume/path")
			assert.False(t, v.Once)
			assert.True(t, ok)
			assert.Contains(t, v.Sources.Content, "*target2_to_be_copied.yaml")
			assert.Contains(t, v.Sources.Content, "*target2_to_be_copied.yml")
			lab, ok := v.Labels["label1"]
			assert.True(t, ok)
			assert.Equal(t, lab, "t2_val1")
			lab, ok = v.Labels["label2"]
			assert.True(t, ok)
			assert.Equal(t, lab, "t2_val2")
		}
	}

	//------------------------------------------------------------
	// Stack1 Hook
	//------------------------------------------------------------
	if assert.Equal(t, 1, len(stack1.Hooks.Deploy.Before)) {
		assert.Equal(t, "task1", stack1.Hooks.Deploy.Before[0].ref)
	}
	if assert.Equal(t, 1, len(stack1.Hooks.Deploy.After)) {
		assert.Equal(t, "task2", stack1.Hooks.Deploy.After[0].ref)
	}

	//------------------------------------------------------------
	// Stack2 Env/Param
	//------------------------------------------------------------
	assert.Equal(t, 2, len(stack2.EnvVars))
	assert.Equal(t, 2, len(stack2.Parameters))
	assert.Equal(t, "stack2_param_key1_value", stack2.Parameters["stack2_param_key1"])
	assert.Equal(t, "stack2_param_key2_value", stack2.Parameters["stack2_param_key2"])
	assert.Equal(t, "stack2_env_key1_value", stack2.EnvVars["stack2_env_key1"])
	assert.Equal(t, "stack2_env_key2_value", stack2.EnvVars["stack2_env_key2"])

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

func buildURL(t *testing.T, loc string) EkURL {
	u, e := CreateUrl(loc)
	assert.Nil(t, e)
	return u
}

func buildRepository(t *testing.T, loc string) Repository {
	base, e := CreateBase("")
	assert.Nil(t, e)
	rep, e := CreateRepository(base, loc, "", "")
	assert.Nil(t, e)
	return rep
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

func TestEnvironmentNameQualifierCustomization(t *testing.T) {
	initial := Environment{Name: "", Qualifier: "", ekara: &Platform{}}
	first := Environment{Name: "FirstName", Qualifier: "FirstQualifier", ekara: &Platform{}}
	second := Environment{Name: "SecondName", Qualifier: "secondQualifier", ekara: &Platform{}}
	initial.Customize(Component{Id: "first"}, &first)
	// The first environment should merge its name and qualifier because those
	// into the initial one are empty.
	assert.Equal(t, "FirstName", initial.Name)
	assert.Equal(t, "FirstQualifier", initial.Qualifier)
	initial.Customize(Component{Id: "second"}, &second)
	// The second environment should NOT merge its name and qualifier because those
	// into the initial one are not empty anymore.
	assert.Equal(t, "FirstName", initial.Name)
	assert.Equal(t, "FirstQualifier", initial.Qualifier)

}
