package pig

//
// Toggles should always be named after the positive feature they represent,
// and default value should be true.
//
type Toggles map[string]bool

func (t Toggles) Any(names ...string) bool {
	for _, name := range names {
		if t[name] {
			return true
		}
	}
	return false
}

func (t Toggles) All(names ...string) bool {
	for _, name := range names {
		if !t[name] {
			return false
		}
	}
	return true
}
