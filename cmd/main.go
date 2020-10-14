package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func main() {
	metadata := logrus.Fields{}
	stateid, _ := metadata["stateid"].(string)
	fmt.Println(len(stateid))
}
