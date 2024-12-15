package storage

import "time"

type Storage interface {
	Shutdown() error
	WriteEmailCode(string, int, time.Time) error
	FindEmailCode(string) (int, time.Time, bool, error)
	DeleteEmailCode(string) error
	WriteNewUser(string, string, string, string) error
	UpdatePassword(string, string) error
	FindUser(string) (string, string, string, bool, error)
}
