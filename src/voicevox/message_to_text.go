package voicevox

import (
	"io/ioutil"
	"log"
)

func msg_to_txt(msg string) {
	byte_msg := []byte(msg)

	err := ioutil.WriteFile("moegi_message.txt", byte_msg, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("File created: moegi_message.txt")
}
