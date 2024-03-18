package loaders

func Get(path string, format string) <-chan string {
	var loaders = map[string]func(string) <-chan string{
		"tranco": CatTranco,
		"txt":    CatTXT,
	}

	if loader, ok := loaders[format]; ok {
		return loader(path)
	}
	return nil
}
