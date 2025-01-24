package models

import "gorm.io/gorm"

// SwiftCode modeluje dane dla kodów SWIFT
type SwiftCode struct {
	gorm.Model
	SwiftCode     string  `gorm:"uniqueIndex;not null"` // Kod SWIFT
	BankName      string  `gorm:"not null"`             // Nazwa banku
	Address       string  // Adres banku
	CountryISO2   string  `gorm:"size:2;not null"` // Kod kraju (ISO2)
	CountryName   string  `gorm:"not null"`        // Nazwa kraju
	TownName      string  // Nazwa miasta
	TimeZone      string  // Strefa czasowa
	IsHeadquarter bool    `gorm:"not null"` // Czy to siedziba główna
	HeadquarterID *string `gorm:"index"`
}
