package definition

type Splits struct {
	Split []Split `xml:"split"`
}

type Split struct {
	Id           int     `xml:"id,attr"`
	Results      Results `xml:"results"`
	PreFunctions struct {
		Functions []Function `xml:"function"`
	} `xml:"pre-functions"`
	PostFunctions struct {
		Functions []Function `xml:"function"`
	} `xml:"post-functions"`
}
