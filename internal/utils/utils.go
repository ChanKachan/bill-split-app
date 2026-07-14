package utils

import (
	"encoding/json"

	entityChat "github.com/ChanKachan/bill-split-app/internal/domain/entity/chat"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ConvertStructsToString(data ...entityChat.Message) []string {
	var values []string
	for _, v := range data {
		jsonData, _ := json.Marshal(v)
		values = append(values, string(jsonData))
	}
	return values
}

func ConvertStringToArrayMessage(data []string) []entityChat.Message {
	var respData []entityChat.Message
	for _, v := range data {
		var msg entityChat.Message

		err := json.Unmarshal([]byte(v), &msg)
		if err != nil {
			return []entityChat.Message{}
		}
		respData = append(respData, msg)
	}

	return respData
}
