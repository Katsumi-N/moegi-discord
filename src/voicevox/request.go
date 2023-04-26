package voicevox

import (
	"net/http"
	"net/url"
	"strconv"
)

const domain = "localhost:50021"

func MakeQuery(text string, speaker int64) error {
	client := &http.Client{}

	query := url.Values{}
	query.Add("speaker", strconv.Itoa(int(speaker)))
	query.Add("text", url.QueryEscape(text))

	req, err := http.NewRequest("POST", domain+"/audio_query", nil)
	if err != nil {
		return err
	}

}
