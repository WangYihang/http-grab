package loaders

import "github.com/WangYihang/gojob/pkg/utils"

func CatTXT(path string) <-chan string {
	return utils.Cat(path)
}
