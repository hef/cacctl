package client

type Option func(client *Client)

func WithUsernameAndPassword(username, password string) Option {
	return func(c *Client) {
		c.username = username
		c.password = password
	}
}
