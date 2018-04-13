package model

type Labels struct {
	labels []string
}

func createLabels(vErrs *ValidationErrors, values ...string) Labels {
	ret := Labels{make([]string, len(values))}
	copy(ret.labels, values)
	return ret
}

func (l Labels) MatchesLabels(candidates ...string) bool {
	for _, l1 := range candidates {
		contains := false
		for _, l2 := range l.labels {
			if l1 == l2 {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}
	return true
}

func (l Labels) AsStrings() []string {
	ret := make([]string, len(l.labels))
	copy(ret, l.labels)
	return ret
}
