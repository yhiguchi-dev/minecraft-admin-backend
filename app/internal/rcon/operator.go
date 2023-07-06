package rcon

import (
	"fmt"
	"strings"
)

type Operator interface {
	GetUserList() ([]string, error)
	GiveItemToUser(userId string, itemId string, amount int) (string, error)
}

type operator struct {
	client Client
}

func NewOperator(client Client) Operator {
	return &operator{client}
}

func (operator *operator) GetUserList() ([]string, error) {
	packet, err := operator.client.Send("list")
	if err != nil {
		return []string{}, err
	}
	payload := string(packet.Payload)
	split := strings.Split(payload, ":")
	_users := strings.Split(split[1], ",")
	users := _users[:]
	for i, v := range users {
		users[i] = strings.TrimSpace(v)
	}
	return users, nil
}

func (operator *operator) GiveItemToUser(userId string, itemId string, amount int) (string, error) {
	packet, err := operator.client.Send(fmt.Sprintf("give %s %s %d", userId, itemId, amount))
	if err != nil {
		return "", err
	}
	payload := string(packet.Payload)
	return payload, nil
}
