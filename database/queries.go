package database

import "time"

func InsertVersionChangeEvent(hash string) {
	query := "INSERT INTO History (hash, createdAt) VALUES (" + hash + "," + time.Now().String() + ")"
	ConnectExecuteAndClose(query)
}
