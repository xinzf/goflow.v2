package definition

type Joins struct {
	Join []Join `xml:"join"`
}

type Join struct {
	Id           int     `xml:"id,attr"`
	Results      Results `xml:"results"`
	PreFunctions struct {
		Functions []Function `xml:"function"`
	} `xml:"pre-functions"`
	PostFunctions struct {
		Functions []Function `xml:"function"`
	} `xml:"post-functions"`
}
