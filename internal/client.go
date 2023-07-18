package internal

type AwsresqClient struct {
}

func NewAwsresqClient() (*AwsresqClient, error) {
	return &AwsresqClient{}, nil
}

func (c *AwsresqClient) Search(service, resource, query string) (string, error) {
	return "", nil
}
