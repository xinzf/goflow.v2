package definition

type Actions struct {
	Actions  []Action `xml:"action"`
	Restrict Restrict `xml:"restrict"`
}

type Action struct {
	Name         string  `xml:"name,attr"`
	Text         string  `xml:"text,attr"`
	Auto         bool    `xml:"auto,attr"`
	Results      Results `xml:"results"`
	PreFunctions struct {
		Functions []Function `xml:"function"`
	} `xml:"pre-functions"`
	PostFunctions struct {
		Functions []Function `xml:"function"`
	} `xml:"post-functions"`
	Restrict Restrict `xml:"restrict"`
}
