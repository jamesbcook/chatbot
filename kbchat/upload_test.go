package kbchat_test

import (
	"os"
	"testing"

	"os/exec"

	"github.com/jamesbcook/chatbot/kbchat"
)

var (
	chatID = os.Getenv("CHATBOT_TEST_CHATID")
)

func TestUpload(t *testing.T) {
	err := exec.Command("wget", "-O", "test_data/golang-250.png", "https://www.idmworks.com/wp-content/themes/idmworks/images/easyblog_images/537/golang-250.png").Run()
	if err != nil {
		t.Fatalf("Error running wget command")
	}
	api, err := kbchat.Start("chat")
	if err != nil {
		t.Fatalf("Failed to start kbchat %v", err)
	}
	fileName := "Testing"
	file := "test_data/golang-250.png"
	if err := api.Upload(chatID, file, fileName); err != nil {
		t.Fatalf("API upload failed with %v", err)
	}

	if err := api.Upload(chatID, file, fileName); err == nil {
		t.Fatalf("API upload failed with %v", err)
	}

	if err := api.Upload("", "shouldn't send", "shouldn't send"); err == nil {
		t.Fatalf("API upload failed with %v", err)
	}

}
