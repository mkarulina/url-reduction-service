package main

func GetLinkByKey(linkKey string) string {
	var foundLink string

	for i := 0; i < len(UrlsDB); i++ {
		if UrlsDB[i].Key == linkKey {
			foundLink = UrlsDB[i].Url
		}
	}
	return foundLink
}
