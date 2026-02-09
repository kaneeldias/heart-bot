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

	message := fmt.Sprintf("%s %s,\n\nThis is your weekly reminder that Kaneel loves you because %s. \n\n%s,\n%s.", salutation, nickname, reason, endPhrase, signature)
	fmt.Println(message)

	recipients := getRecipients()
	bcc := getBCC()

	fmt.Printf("Recipients: %s\n", recipients)
	fmt.Printf("BCC: %s\n\n", bcc)

	SendEmail(subject, message, recipients, bcc)
}
