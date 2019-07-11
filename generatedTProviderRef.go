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
// Implementation(s) of TProviderRef
// ----------------------------------------------------

//TProviderRefOnProviderRefHolder is the struct containing the ProviderRef in order to implement TProviderRef
type TProviderRefOnProviderRefHolder struct {
	h ProviderRef
}

//CreateTProviderRefForProviderRef returns an holder of ProviderRef implementing TProviderRef
func CreateTProviderRefForProviderRef(o ProviderRef) TProviderRefOnProviderRefHolder {
	return TProviderRefOnProviderRefHolder{
		h: o,
	}
}

//Provider returns the provider wherein the node should be deployed
func (r TProviderRefOnProviderRefHolder) Provider() (TProvider, error) {
	v, err := r.h.Resolve()
	return CreateTProviderForProvider(v), err
}