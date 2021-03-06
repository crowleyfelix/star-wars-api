package controllers

//RequestContext exposes request context methods
type RequestContext interface {
	NegotiateFormat(...string) string
	ShouldBindJSON(interface{}) error
	JSON(int, interface{})
	Header(string, string)
	Data(int, string, []byte)
	DefaultQuery(string, string) string
	BindQuery(interface{}) error
	BindJSON(interface{}) error
	Param(string) string
}
