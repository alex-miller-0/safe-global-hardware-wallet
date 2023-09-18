package util

func PrintTabs(tabs int) string {
	str := ""
	for i := 0; i < tabs; i++ {
		str += "  "
	}
	return str
}
