package main

import "fmt"

type Device struct {
	Id string
}

func GetDevices() map[string]Device {
	count := 10
	result := map[string]Device{}

	for i := 0; i < count; i++ {
		id := fmt.Sprintf("device_id_%v", i)
		result[id] = Device{
			Id: id,
		}
	}

	return result
}
