package service

type ResultList struct {
	Service string `json:"service"`
	Resource string `json:"resource"`
	Results []interface{} `json:"results"`
}

type AwsAPI interface {
	Validate(resource string) bool
	Query(resource string) (*ResultList, error)
}
