package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateReferencesComplete(t *testing.T) {
	envRefs, e := ParseYamlDescriptorReferences(buildRepository(t, "./testdata/yaml/complete.yaml").Url, &TemplateContext{})
	assert.Nil(t, e)
	assertEnvRefs(t, envRefs)
}

func assertEnvRefs(t *testing.T, env EnvironmentReferences) {

	// Platform
	assert.NotNil(t, env.Ekara)
	assert.NotNil(t, env.Ekara.Components)
	assert.Equal(t, 8, len(env.Ekara.Components))
	assert.Equal(t, "someBase", env.Ekara.Base)
	assert.Equal(t, "ekara-platform/distribution", env.Ekara.Parent.Repository)

	// Variables
	assert.NotNil(t, env.Vars)
	if assert.Len(t, env.Vars, 2) {
		va, ok := env.Vars["global_var_key1"]
		assert.True(t, ok)
		assert.Equal(t, va, "global_var_val1")

		va, ok = env.Vars["global_var_key2"]
		assert.True(t, ok)
		assert.Equal(t, va, "global_var_val2")
	}

	//------------------------------------------------------------
	// Orchestrator
	//------------------------------------------------------------
	orchestrator := env.OrchestratorRefs
	assert.NotNil(t, orchestrator)

	//------------------------------------------------------------
	// Environment Providers
	//------------------------------------------------------------
	providers := env.ProvidersRefs
	assert.NotNil(t, providers)
	assert.Len(t, providers, 2)

	assert.Contains(t, providers, "aws", "azure")
	assert.Equal(t, "aws", providers["aws"].Component)
	assert.Equal(t, "azure", providers["azure"].Component)

	//------------------------------------------------------------
	// Nodes
	//------------------------------------------------------------
	nodes := env.NodesRefs
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 2)

	assert.Contains(t, nodes, "node1", "node2")
	assert.Equal(t, "aws", nodes["node1"].Provider.Component)
	assert.Equal(t, "azure", nodes["node2"].Provider.Component)

	//------------------------------------------------------------
	// Stacks
	//------------------------------------------------------------
	stacks := env.StacksRefs
	assert.NotNil(t, stacks)
	assert.Len(t, stacks, 2)

	assert.Contains(t, stacks, "stack1", "stack2")
	assert.Equal(t, "stack1", stacks["stack1"].Component)
	assert.Equal(t, "stack2", stacks["stack2"].Component)

	//------------------------------------------------------------
	// Tasks
	//------------------------------------------------------------
	tasks := env.TasksRefs
	assert.NotNil(t, tasks)
	assert.Len(t, tasks, 3)
	assert.Contains(t, tasks, "task1", "task2", "task3")
	assert.Equal(t, "task1", tasks["task1"].Component)
	assert.Equal(t, "task2", tasks["task2"].Component)
	assert.Equal(t, "task3", tasks["task3"].Component)

	//------------------------------------------------------------
	// Used Referencies
	//------------------------------------------------------------
	used, orphans := env.Uses(&Orphans{})
	assert.Len(t, orphans.Refs, 0)
	assert.Len(t, used.Refs, 8)
	assert.Contains(t, used.Refs, "aws", "azure", "swarm", "stack1", "stack2", "task1", "task2", "task3")
}
