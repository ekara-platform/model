package model

const (
	//DefaultComponentBase specifies the default base URL where to look for a component
	//
	// For example if compoment is defined like this:
	//  components:
	//    aws:
	//      repository: ekara-platform/aws-provider
	//      version: 1.2.3
	//
	// We will assume that this is a Git compoment located in:
	//   https://github.com/: ekara-platform/aws-provider
	//
	DefaultComponentBase = "https://github.com"

	//CoreComponentId identifies the main component, base of the ekara platform
	CoreComponentId = "__core__"

	//CoreComponentRepo identifies the repository where to fetch the main component of the ekara platform
	CoreComponentRepo = "ekara-platform/core"

	//DefaultDescriptorName specifies the default name of the environment descriptor
	//
	//When the environment descriptor is not specified, for example into a use
	// component then we will look for a default descriptor name "ekara.yaml"
	DefaultDescriptorName = "ekara.yaml"

	//GitExtentsion represents the extension of the GIT repository extension
	GitExtentsion = ".git"
)
