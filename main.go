package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Variables struct {
	Subjects    []string `json:"subjects"`
	Salutations []string `json:"salutations"`
	Nicknames   []string `json:"nicknames"`
	Reasons     []string `json:"reasons"`
	EndPhrases  []string `json:"end_phrases"`
	Signatures  []string `json:"signatures"`
}

func main() {
	file, err := os.Open("variables.json")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	var variables Variables
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&variables)
	if err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	subject := fmt.Sprintf("=?utf-8?q?%s_=E2=9D=A4?=", GetRandomFromList(variables.Subjects))
	salutation := GetRandomFromList(variables.Salutations)
	nickname := GetRandomFromList(variables.Nicknames)
	reason := GetRandomFromList(variables.Reasons)
	endPhrase := GetRandomFromList(variables.EndPhrases)
	signature := GetRandomFromList(variables.Signatures)

	htmlBody := fmt.Sprintf("<p>%s %s,</p><p>This is your weekly reminder that Kaneel loves you because %s.</p><p>%s,<br>%s.</p> <img src=\"cid:image1\" alt=\"Kaneel loves you\" />", salutation, nickname, reason, endPhrase, signature)
	fmt.Println(htmlBody)

	imageGenDescription := fmt.Sprintf("Create a picture of a boy and a girl %s in a pixar style art", reason)
	GenerateImage(imageGenDescription)

	base64Image, err := getFileInBase64("generated_image.png")
	if err != nil {
		fmt.Printf("Error reading image file: %v\n", err)
		base64Image = ""
	}

	recipients := getRecipients()
	bcc := getBCC()

	fmt.Printf("Recipients: %s\n", recipients)
	fmt.Printf("BCC: %s\n\n", bcc)

	SendEmail(subject, htmlBody, recipients, bcc, base64Image)
}

func getFileInBase64(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
