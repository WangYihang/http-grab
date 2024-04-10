package loader

func CatTranco(path string) <-chan string {
	return CatCSV(path, 1)
}
