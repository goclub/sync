package context_test

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestWithTimeout(t *testing.T) {
	{
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		s, err := Call(ctx, Option{
			Name: "1",
			ReturnError: false,
			Sleep: time.Second*1,
		}) ; if err != nil {
			log.Print("1 error:", err)
		} else {
			log.Print("1 result:", s)
		}
	}
	{
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		s, err := Call(ctx, Option{
			Name: "2",
			ReturnError: false,
			Sleep: time.Second*3,
		}) ; if err != nil {
		log.Print("2 error:", err)
		} else {
			log.Print("2 result:", s)
		}
	}
}