package descriptor

type Lagoon struct {
	root *Environment
	Component

	Proxy Proxy
}

type Proxy struct {
	Http    string
	Https   string
	NoProxy string
}

func createLagoon(vErrs *ValidationErrors, env *Environment, yamlEnv *yamlEnvironment) Lagoon {
	return Lagoon{
		root:      env,
		Component: createComponent(vErrs, "lagoon-platform", createVersion(vErrs, "lagoon", yamlEnv.Version)),
		Proxy:     Proxy{Http: yamlEnv.Lagoon.Proxy.Http, Https: yamlEnv.Lagoon.Proxy.Https, NoProxy: yamlEnv.Lagoon.Proxy.NoProxy}}
}
