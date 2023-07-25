package api

var client *Client

func init() {
	var err error
	client, err = NewFromEnv()
	if err != nil {
		panic(err)
	}
}
