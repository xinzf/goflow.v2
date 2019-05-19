package definition

type Restrict struct {
	Message    string     `xml:"message,attr"`
	Conditions Conditions `xml:"conditions"`
}

func (static Restrict) GetMessage() string {
	if static.Message == "" {
		return "Restrict verify failed."
	}

	return static.Message
}
