package definition

type Global struct {
	PreFunctions struct {
		Functions []Function `xml:"function"`
	} `xml:"pre-functions"`
	PostFunctions struct {
		Functions []Function `xml:"function"`
	} `xml:"post-functions"`
}
