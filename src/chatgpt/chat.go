package chatgpt

import (
	"bufio"
	"log"
	"os"
)

func Chat(msgs []string) ([]string, error) {
	file, err := os.Open("moegi-settings.txt")
	if err != nil {
		log.Println("can't read setting.txt")
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var setting string
	for scanner.Scan() {
		setting += scanner.Text()
	}
	send := []SendingMessage{{Role: "system", Content: setting}}

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
