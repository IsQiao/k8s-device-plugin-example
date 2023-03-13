package main

import (
	"fmt"
	"os"
)

type Device struct {
	Id string
}

func GetDevices() map[string]Device {
	count := 10
	result := map[string]Device{}

	for i := 0; i < count; i++ {
		id := fmt.Sprintf("%s_%v", getHostName(), i)
		result[id] = Device{
			Id: id,
		}
	}

	return result
}

func getHostName() string {
	hostname, _ := os.Hostname()
	return hostname
}
