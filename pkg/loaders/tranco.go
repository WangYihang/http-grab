package loaders

func CatTranco(path string) <-chan string {
	return CatCSV(path, 1)
}
