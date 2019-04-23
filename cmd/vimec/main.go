package main

import "log"

func main(){
	if err := runMainWindow(); err != nil{
		log.Fatal(err)
	}
}