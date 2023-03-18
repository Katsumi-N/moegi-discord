package chatgpt

import (
	"io/ioutil"
	"log"
)

func Chat(msgs []string) ([]string, error) {
	setting, err := ioutil.ReadFile("moegi-settings.txt")
	if err != nil {
		log.Println("can't read setting.txt")
		return []string{}, err
	}

	send := make([]SendingMessage, 0)
	send = append(send, SendingMessage{
		Role:    "system",
		Content: string(setting),
	})
	for _, m := range msgs {
		send = append(send, SendingMessage{
			Role:    "user",
			Content: m,
		})
	}

	resp, err := RequestToOpenAI(send)
	if err != nil {
		log.Fatal(err)
		return []string{}, err
	}

	var resp_msgs []string
	for _, m := range resp.Choices {
		resp_msgs = append(resp_msgs, m.Message.Content)
	}
	return resp_msgs, err
}
