package definition

type Arg struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}
