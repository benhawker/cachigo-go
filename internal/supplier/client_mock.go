package supplier

// MockClient ...
type MockClient struct {
	MakeRequestFunc func(string) (Response, error)
}

// MakeRequest ...
func (c MockClient) MakeRequest(url string) (Response, error) {
	return c.MakeRequestFunc(url)
}
