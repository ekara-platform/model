package model

//*****************************************************************************
//
//     _         _           ____                           _           _
//    / \  _   _| |_ ___    / ___| ___ _ __   ___ _ __ __ _| |_ ___  __| |
//   / _ \| | | | __/ _ \  | |  _ / _ \ '_ \ / _ \ '__/ _` | __/ _ \/ _` |
//  / ___ \ |_| | || (_) | | |_| |  __/ | | |  __/ | | (_| | ||  __/ (_| |
// /_/   \_\__,_|\__\___/   \____|\___|_| |_|\___|_|  \__,_|\__\___|\__,_|
//
// This file is autogenerated by "go generate .". Do not modify.
//
//*****************************************************************************

// ----------------------------------------------------
// Implementation(s) of TPlatform
// ----------------------------------------------------

//TPlatformOnPlatformHolder is the struct containing the Platform in order to implement TPlatform
type TPlatformOnPlatformHolder struct {
	h Platform
}

//CreateTPlatformForPlatform returns an holder of Platform implementing TPlatform
func CreateTPlatformForPlatform(o Platform) TPlatformOnPlatformHolder {
	return TPlatformOnPlatformHolder{
		h: o,
	}
}

//Base returns the base location of the platform
func (r TPlatformOnPlatformHolder) Base() TBase {
	return CreateTBaseForBase(r.h.Base)
}

//Parent returns the parent used by the platform
func (r TPlatformOnPlatformHolder) Parent() TComponent {
	return CreateTComponentForParent(r.h.Parent)
}

//HasComponents returns true if the platform has components
func (r TPlatformOnPlatformHolder) HasComponents() bool {
	return len(r.h.Components) > 0
}

//HasTemplates returns true if the environment has defined templates
func (r TPlatformOnPlatformHolder) HasTemplates() bool {
	return len(r.h.Templates) > 0
}

//Templates returns the environment templates
func (r TPlatformOnPlatformHolder) Templates() []string {
	return r.h.Templates
}
