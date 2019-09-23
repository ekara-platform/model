package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var stackOtherT1, stackOtherT2, stackOriginT1, stackOriginT2 TaskRef

func TestStackDescType(t *testing.T) {
	s := Stack{}
	assert.Equal(t, s.DescType(), "Stack")
}

func TestStackDescName(t *testing.T) {
	s := Stack{Name: "my_name"}
	assert.Equal(t, s.DescName(), "my_name")
}

func getStackOrigin() Stack {
	s1 := Stack{Name: "my_name"}
	s1.EnvVars = make(map[string]string)
	s1.EnvVars["key1"] = "val1_target"
	s1.EnvVars["key2"] = "val2_target"
	s1.Parameters = make(map[string]interface{})
	s1.Parameters["key1"] = "val1_target"
	s1.Parameters["key2"] = "val2_target"

	s1.cRef = componentRef{
		ref: "cOriginal",
	}
	depS1 := []string{"dep1", "dep2"}
	s1.DependsOn = createDependencies(&Environment{}, DescriptorLocation{Path: "location"}, "stack_name", depS1)

	stackOriginT1 = TaskRef{ref: "T1"}
	stackOriginT2 = TaskRef{ref: "T2"}
	hS1 := StackHook{}
	hS1.Deploy.Before = append(hS1.Deploy.Before, stackOriginT1)
	hS1.Deploy.After = append(hS1.Deploy.After, stackOriginT2)
	s1.Hooks = hS1

	s1.Copies = Copies{
		Content: make(map[string]Copy),
	}
	s1.Copies.Content["path1"] = Copy{
		Labels: map[string]string{
			"label_key1": "label_value1",
			"label_key2": "label_value2",
		},
		Sources: Patterns{
			Content: []string{"path1", "path2"},
		},
	}
	return s1
}

func getStackOther(name string) Stack {
	other := Stack{Name: name}
	other.EnvVars = make(map[string]string)
	other.EnvVars["key2"] = "val2_other"
	other.EnvVars["key3"] = "val3_other"
	other.Parameters = make(map[string]interface{})
	other.Parameters["key2"] = "val2_other"
	other.Parameters["key3"] = "val3_other"
	other.cRef = componentRef{
		ref: "cOther",
	}
	depOther := []string{"dep2", "dep3"}
	other.DependsOn = createDependencies(&Environment{}, DescriptorLocation{Path: "location"}, "stack_name", depOther)

	stackOtherT1 = TaskRef{ref: "T3"}
	stackOtherT1 = TaskRef{ref: "T4"}
	oS1 := StackHook{}
	oS1.Deploy.Before = append(oS1.Deploy.Before, stackOtherT1)
	oS1.Deploy.After = append(oS1.Deploy.After, stackOtherT2)
	other.Hooks = oS1

	other.Copies = Copies{
		Content: make(map[string]Copy),
	}
	other.Copies.Content["path2"] = Copy{
		Labels: map[string]string{
			"label_key3": "label_value3",
			"label_key4": "label_value4",
		},
		Sources: Patterns{
			Content: []string{"path3", "path4"},
		},
	}
	return other
}

func checkStackMerge(t *testing.T, s Stack) {
	assert.Equal(t, s.cRef.ref, "cOther")
	if assert.Len(t, s.DependsOn.Content, 3) {
		assert.Equal(t, s.DependsOn.Content[0].ref, "dep1")
		assert.Equal(t, s.DependsOn.Content[1].ref, "dep2")
		assert.Equal(t, s.DependsOn.Content[2].ref, "dep3")
	}

	if assert.Len(t, s.EnvVars, 3) {
		checkMap(t, s.EnvVars, "key1", "val1_target")
		checkMap(t, s.EnvVars, "key2", "val2_other")
		checkMap(t, s.EnvVars, "key3", "val3_other")
	}

	if assert.Len(t, s.Parameters, 3) {
		checkMapInterface(t, s.Parameters, "key1", "val1_target")
		checkMapInterface(t, s.Parameters, "key2", "val2_other")
		checkMapInterface(t, s.Parameters, "key3", "val3_other")
	}

	if assert.Len(t, s.Hooks.Deploy.Before, 2) {
		assert.Contains(t, s.Hooks.Deploy.Before, stackOriginT1, stackOtherT1)
		assert.Equal(t, s.Hooks.Deploy.Before[0], stackOriginT1)
		assert.Equal(t, s.Hooks.Deploy.Before[1], stackOtherT1)
	}

	if assert.Len(t, s.Hooks.Deploy.After, 2) {
		assert.Contains(t, s.Hooks.Deploy.After, stackOriginT2, stackOtherT2)
		assert.Equal(t, s.Hooks.Deploy.After[0], stackOriginT2)
		assert.Equal(t, s.Hooks.Deploy.After[1], stackOtherT2)
	}

	assert.Len(t, s.Copies.Content, 2)
}

func TestStackMerge(t *testing.T) {
	o := getStackOrigin()
	err := o.customize(getStackOther("my_name"))
	if assert.Nil(t, err) {
		checkStackMerge(t, o)
	}
}

func TestMergeStackItself(t *testing.T) {
	o := getStackOrigin()
	oi := o
	err := o.customize(o)
	if assert.Nil(t, err) {
		assert.Equal(t, oi, o)
	}
}

func TestStacksMerge(t *testing.T) {
	origins := make(Stacks)
	origins["myS"] = getStackOrigin()
	others := make(Stacks)
	others["myS"] = getStackOther("my_name")

	customized, err := origins.customize(&Environment{}, others)
	if assert.Nil(t, err) {
		if assert.Len(t, customized, 1) {
			o := customized["myS"]
			checkStackMerge(t, o)
		}
	}
}

func TestStacksMergeAddition(t *testing.T) {
	origins := make(Stacks)
	origins["myS"] = getStackOrigin()
	others := make(Stacks)
	others["myS"] = getStackOther("my_name")
	others["new"] = getStackOther("new")

	customized, err := origins.customize(&Environment{}, others)
	if assert.Nil(t, err) {
		assert.Len(t, customized, 2)
	}
}

func TestStacksEmptyMerge(t *testing.T) {
	origins := make(Stacks)
	o := getStackOrigin()
	origins["myS"] = o
	others := make(Stacks)

	customized, err := origins.customize(&Environment{}, others)
	if assert.Nil(t, err) {
		if assert.Len(t, customized, 1) {
			oc := customized["myS"]
			assert.Equal(t, o, oc)
		}
	}
}

func TestMergeStackUnrelated(t *testing.T) {
	sta := Stack{
		Name: "Name",
	}

	o := Stack{
		Name: "Dummy",
	}

	err := sta.customize(o)
	if assert.NotNil(t, err) {
		assert.Equal(t, err.Error(), "cannot customize unrelated stacks (Name != Dummy)")
	}
}

func TestStacksDependencies(t *testing.T) {
	sts := Stacks{}

	// This test will check the dependencies on the following tree
	//
	//        0                    1
	//                           / | \
	//                          2  7  8
	//                        / |     | \
	//                       3  6     9  12
	//                      / |       | \
	//                     4  5       10  11

	sts["6"] = Stack{Name: "6", DependsOn: getDependsOn("2")}
	sts["1"] = Stack{Name: "1"}
	sts["3"] = Stack{Name: "3", DependsOn: getDependsOn("2")}
	sts["4"] = Stack{Name: "4", DependsOn: getDependsOn("3")}
	sts["9"] = Stack{Name: "9", DependsOn: getDependsOn("8")}
	sts["5"] = Stack{Name: "5", DependsOn: getDependsOn("3")}
	sts["7"] = Stack{Name: "7", DependsOn: getDependsOn("1")}
	sts["8"] = Stack{Name: "8", DependsOn: getDependsOn("1")}
	sts["10"] = Stack{Name: "10", DependsOn: getDependsOn("9")}
	sts["2"] = Stack{Name: "2", DependsOn: getDependsOn("1")}
	sts["12"] = Stack{Name: "12", DependsOn: getDependsOn("8")}
	sts["11"] = Stack{Name: "11", DependsOn: getDependsOn("9")}
	sts["0"] = Stack{Name: "0"}

	assert.Len(t, sts, 13)
	resolved, err := sts.ResolveDependencies()
	assert.Equal(t, len(resolved), len(sts))
	assert.Nil(t, err)

	processed := make(map[string]Stack)
	for _, val := range resolved {
		if len(processed) == 0 {
			processed[val.Name] = val
			continue
		}
		//Check than all the dependencies has been already processd
		for _, d := range val.DependsOn.Content {
			if _, ok := processed[d.ref]; !ok {
				assert.Fail(t, "Dependency %s has not been yet processed", d.ref)
			}

		}
		processed[val.Name] = val
	}
	//Check that the original Stacks has been untouched
	assert.Len(t, sts, 13)

}

func TestStacksMultiplesDependencies(t *testing.T) {
	sts := Stacks{}

	// This test will check the dependencies on the following tree
	//
	//        0                    1
	//                           / | \
	//                          2  7  8
	//                        / | /    | \
	//                       3  6     9  12
	//                      / |   \   | \
	//                     4  5    \-10  11

	sts["6"] = Stack{Name: "6", DependsOn: getDependsOn("2", "7")}
	sts["1"] = Stack{Name: "1"}
	sts["3"] = Stack{Name: "3", DependsOn: getDependsOn("2")}
	sts["4"] = Stack{Name: "4", DependsOn: getDependsOn("3")}
	sts["9"] = Stack{Name: "9", DependsOn: getDependsOn("8")}
	sts["5"] = Stack{Name: "5", DependsOn: getDependsOn("3")}
	sts["7"] = Stack{Name: "7", DependsOn: getDependsOn("1")}
	sts["8"] = Stack{Name: "8", DependsOn: getDependsOn("1")}
	sts["10"] = Stack{Name: "10", DependsOn: getDependsOn("9", "6")}
	sts["2"] = Stack{Name: "2", DependsOn: getDependsOn("1")}
	sts["12"] = Stack{Name: "12", DependsOn: getDependsOn("8")}
	sts["11"] = Stack{Name: "11", DependsOn: getDependsOn("9")}
	sts["0"] = Stack{Name: "0"}

	assert.Len(t, sts, 13)
	resolved, err := sts.ResolveDependencies()
	assert.Equal(t, len(resolved), len(sts))
	assert.Nil(t, err)

	processed := make(map[string]Stack)
	for _, val := range resolved {
		if len(processed) == 0 {
			processed[val.Name] = val
			continue
		}
		//Check than all the dependencies has been already processd
		for _, d := range val.DependsOn.Content {
			if _, ok := processed[d.ref]; !ok {
				assert.Fail(t, "Dependency %s has not been yet processed", d.ref)
			}

		}
		processed[val.Name] = val
	}
	//Check that the original Stacks has been untouched
	assert.Len(t, sts, 13)

}

func TestStacksCyclicDependencies(t *testing.T) {
	sts := Stacks{}

	sts["1"] = Stack{Name: "1", DependsOn: getDependsOn("3")}
	sts["2"] = Stack{Name: "2", DependsOn: getDependsOn("1")}
	sts["3"] = Stack{Name: "3", DependsOn: getDependsOn("2")}

	assert.Len(t, sts, 3)
	_, err := sts.ResolveDependencies()
	assert.NotNil(t, err)

	//Check that the original Stacks has been untouched
	assert.Len(t, sts, 3)

}

func getDependsOn(deps ...string) Dependencies {
	res := make([]StackRef, 0)
	for _, s := range deps {
		r := StackRef{
			ref: s,
		}
		res = append(res, r)
	}
	return Dependencies{Content: res}
}
