package model

// A HumanDescriber is a type that can describe itself as a string which can
// be easily understood or interpreted by a human.
type HumanDescriber interface {
	HumanDescribe() string
}
