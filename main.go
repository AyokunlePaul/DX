package main

import (
	"DX/src/api/application"
)

func main() {
	//viper.SetConfigFile(".env")
	//if err := viper.ReadInConfig(); err != nil {
	//	panic(err)
	//}
	application.StartApplication()
}
