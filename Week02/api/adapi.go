package api

import (
	"Week02/biz"
	"fmt"
)

func Ad() error {
	lis, err := biz.Advertisement(3)
	if err != nil {
		return err
	}
	fmt.Println(lis)
	return nil
}


