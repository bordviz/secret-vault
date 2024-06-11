package main

import (
	"fmt"
	"os"
	"vault/internal/config"
)

func main() {
	f, err := os.Create("cmd/migrations/up.sh")
	if err != nil {
		fmt.Println(err)
		return
	}

	cfg := config.MustLoad()

	upSh :=
		fmt.Sprintf(
			`#!/bin/bash 

goose -dir migrations postgres "postgresql://%s:%s@%s:%d/%s?sslmode=%s" up
			`,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		)

	l, err := f.WriteString(upSh)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err = os.Create("cmd/migrations/down.sh")
	if err != nil {
		fmt.Println(err)
		return
	}

	upDown :=
		fmt.Sprintf(
			`#!/bin/bash 

goose -dir migrations postgres "postgresql://%s:%s@%s:%d/%s?sslmode=%s" down
			`,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		)

	l, err = f.WriteString(upDown)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
