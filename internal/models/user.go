package models

import "time"

type User struct {
	ID    int       `json:"id"`
	Name  string    `json:"name"`
	Age   int       `json:"age"`
	City  string    `json:"city"`
	RegDt time.Time `json:"reg_dt"`
}
