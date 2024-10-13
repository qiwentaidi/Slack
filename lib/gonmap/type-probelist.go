package gonmap

type ProbeList []string

func (p ProbeList) exist(probeName string) bool {
	for _, name := range p {
		if name == probeName {
			return true
		}
	}
	return false
}
