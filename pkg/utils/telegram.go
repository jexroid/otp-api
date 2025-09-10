package utils

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func Telegram(message string) {
	botToken := os.Getenv("BOT_TOKEN")
	chatId := os.Getenv("CHAT_ID")
	var url = `https://api.telegram.org/bot` + botToken + `/sendMessage?chat_id=` + chatId + `&text=` + message
	content := `Url: ` + url + `, AgentList: "Mozilla Firefox", VersionList: "HTTP/1.1", MethodList: "GET"`
	jsonStr := []byte(content)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Panic(err)
	}
	defer resp.Body.Close()
}
