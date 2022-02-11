package main

func getLinkById(linkId string) string {
	var foundLink string

	for i := 0; i < len(urlsDB); i++ {
		if urlsDB[i].Id == linkId {
			foundLink = urlsDB[i].Url
		}
	}
	return foundLink
}
