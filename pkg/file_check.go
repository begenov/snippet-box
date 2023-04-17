package pkg

import (
	"fmt"
	"os"
)

const filename = "../log_file.log"

func Read() {
	read, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	fmt.Print(string(read))
}

func Removefile() {
	err := os.Remove(filename)
	if err != nil {
		fmt.Println("error: ok")
		return
	}
}
