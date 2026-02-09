package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func getClient(config *oauth2.Config) *http.Client {
	//tokFile := "token.json"
	//tok, err := tokenFromFile(tokFile)
	//if err != nil {
	//	tok = getTokenFromWeb(config)
	//	saveToken(tokFile, tok)
	//}

	token := &oauth2.Token{
		RefreshToken: os.Getenv("GOOGLE_REFRESH_TOKEN"),
	}
	return config.Client(context.Background(), token)
}

//func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
//	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
//	fmt.Printf("Go to the following link in your browser then type the "+
//		"authorization code: \n%v\n", authURL)
//
//	var authCode string
//	if _, err := fmt.Scan(&authCode); err != nil {
//		log.Fatalf("Unable to read authorization code: %v", err)
//	}
//
//	tok, err := config.Exchange(context.TODO(), authCode)
//	if err != nil {
//		log.Fatalf("Unable to retrieve token from web: %v", err)
//	}
//	return tok
//}
//
//func tokenFromFile(file string) (*oauth2.Token, error) {
//	f, err := os.Open(file)
//	if err != nil {
//		return nil, err
//	}
//	defer f.Close()
//	tok := &oauth2.Token{}
//	err = json.NewDecoder(f).Decode(tok)
//	return tok, err
//}
//
//func saveToken(path string, token *oauth2.Token) {
//	fmt.Printf("Saving credential file to: %s\n", path)
//	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
//	if err != nil {
//		log.Fatalf("Unable to cache oauth token: %v", err)
//	}
//	defer f.Close()
//	json.NewEncoder(f).Encode(token)
//}

func SendEmail(subject string, htmlBody string, recipients []string, bcc []string, base64Image string) {
	ctx := context.Background()

	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{gmail.GmailSendScope},
	}

	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	message := makeRawMessageWithInlineImage(subject, htmlBody, recipients, bcc, "generated_image.png")

	_, err = srv.Users.Messages.Send("me", message).Do()
	if err != nil {
		log.Fatalf("Unable to send email: %v", err)
	}
}

func makeRawMessageWithInlineImage(subject, htmlBody string, recipients, bcc []string, imageFilePath string) *gmail.Message {
	to := strings.Join(recipients, ", ")
	bccHeader := strings.Join(bcc, ", ")

	// Read image file
	data, err := os.ReadFile(imageFilePath)
	if err != nil {
		log.Printf("Unable to read image file: %v", err)
		data = []byte("")
	}

	// Detect MIME type from the first up to 512 bytes
	head := data
	if len(head) > 512 {
		head = head[:512]
	}
	mimeType := http.DetectContentType(head)

	// Extract filename without adding new imports
	filename := imageFilePath
	if idx := strings.LastIndexAny(imageFilePath, "/\\"); idx != -1 {
		filename = imageFilePath[idx+1:]
	}
	if filename == "" {
		filename = "image"
	}

	// Base64 encode the image
	b64Image := base64.StdEncoding.EncodeToString(data)

	boundary := "BOUNDARY_STRING"

	var sb strings.Builder
	// Headers + multipart/related
	sb.WriteString(fmt.Sprintf("To: %s\r\nBcc: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/related; boundary=%s\r\n\r\n", to, bccHeader, subject, boundary))

	// HTML part (refer to image via cid:image1)
	sb.WriteString(fmt.Sprintf("--%s\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s\r\n\r\n", boundary, htmlBody))

	// Image part (inline, base64)
	sb.WriteString(fmt.Sprintf("--%s\r\nContent-Type: %s; name=\"%s\"\r\nContent-Transfer-Encoding: base64\r\nContent-ID: <image1>\r\nContent-Disposition: inline; filename=\"%s\"\r\n\r\n%s\r\n", boundary, mimeType, filename, filename, b64Image))

	// end
	sb.WriteString(fmt.Sprintf("--%s--", boundary))

	raw := base64.URLEncoding.EncodeToString([]byte(sb.String()))
	raw = strings.TrimRight(raw, "=") // Gmail expects URL-safe base64 without padding

	return &gmail.Message{Raw: raw}
}
