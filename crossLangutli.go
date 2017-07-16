package libv2ray

type StringArrayList struct {
	list []string
}

func (al *StringArrayList) GetLen() int {
	return len(al.list)
}

func (al *StringArrayList) GetElementById(id int) string {
	return al.list[id]
}
