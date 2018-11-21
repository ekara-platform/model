package model

type (
	Reference struct {
		Id        string
		Type      string
		Mandatory bool
		Location  DescriptorLocation
		Repo      map[string]interface{}
	}

	ValidReference interface {
		Reference() Reference
	}
)
