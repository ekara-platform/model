package descriptor

type Labels struct {
	labels []string
}

func createLabels(values ...string) Labels {
	ret := Labels{make([]string, len(values))}
	copy(ret.labels, values)
	return ret
}

func (l Labels) Contains(candidates ...string) bool {
	for _, l1 := range candidates {
		contains := false
		for _, l2 := range l.labels {
			if l1 == l2 {
				contains = true
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
