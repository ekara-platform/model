package descriptor

type Lagoon struct {
	root *Environment
	Component

	Proxy struct {
		Http    string
		Https   string
		NoProxy string
	}
}

func createLagoon(env *Environment, yamlEnv *yamlEnvironment) (res Lagoon, err error) {
	res = Lagoon{root: env}

	res.Version, err = createVersion(yamlEnv.Lagoon.Version)
	if err != nil {
		return
	}

	res.Proxy.Http = yamlEnv.Lagoon.Proxy.Http
	res.Proxy.Https = yamlEnv.Lagoon.Proxy.Https
	res.Proxy.NoProxy = yamlEnv.Lagoon.Proxy.NoProxy

	return
}
