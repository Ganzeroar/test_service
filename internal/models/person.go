package models

import "gorm.io/gorm"

type Person struct {
	gorm.Model
	Name        string
	Surname     string
	Patronymic  string
	Age         int8
	Gender      string
	Nationality string
}

func (Person) TableName() string {
	return "persons"
}
