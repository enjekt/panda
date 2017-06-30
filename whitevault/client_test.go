package main

import (
	"testing"
	"github.com/enjekt/commons"
	"fmt"
)

func TestSendTokenAndPad(t *testing.T) {
	token:= commons.InitToken("2222222222")
	pad := commons.InitPad("11111111")

	SendTokenAndPad(token,pad)
}

func TestGetPad(t *testing.T) {
	token:= commons.InitToken("2222222222")
	pad:=GetPad(token)
	fmt.Println(pad.ToString())
}