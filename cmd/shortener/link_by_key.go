package main

func (c *Container) GetLinkByKey(linkKey string) string {
	var foundLink string

	if val, found := c.urls[linkKey]; found {
		foundLink = val
	}
	return foundLink
}
