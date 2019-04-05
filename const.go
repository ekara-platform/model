package model

const (
	//DefaultComponentBase specifies the default base URL where to look for a component
	//
	// For example if component is defined like this:
	//  components:
	//    aws:
	//      repository: ekara-platform/aws-provider
	//      version: 1.2.3
	//
	// We will assume that this is a Git component located in:
	//   https://github.com/: ekara-platform/aws-provider
	//
	DefaultComponentBase = "https://github.com"

	//DefaultDescriptorName specifies the default name of the environment descriptor
	//
	//When the environment descriptor is not specified, for example into a use
	// component then we will look for a default descriptor name "ekara.yaml"
	DefaultDescriptorName = "ekara.yaml"

	//EkaraComponentId The component identifier for the ekara distribution
	EkaraComponentId = "__ekara__"

	//EkaraComponentRepo The default repository for the ekara distribution
	EkaraComponentRepo = "ekara-platform/distribution"

	//GitExtension represents the extension of the GIT repository extension
	GitExtension = ".git"
)
