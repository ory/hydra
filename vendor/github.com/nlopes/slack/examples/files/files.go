package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

func main() {
	api := slack.New("YOUR_TOKEN_HERE")
	params := slack.FileUploadParameters{
		Title: "Batman Example",
		//Filetype: "txt",
		File: "example.txt",
		//Content:  "Nan Nan Nan Nan Nan Nan Nan Nan Batman",
	}
	file, err := api.UploadFile(params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Name: %s, URL: %s\n", file.Name, file.URL)

	err = api.DeleteFile(file.ID)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("File %s deleted successfully.\n", file.Name)
}
