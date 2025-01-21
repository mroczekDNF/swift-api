package models

import "gorm.io/gorm"

// SwiftCode modeluje dane dla kodów SWIFT
type SwiftCode struct {
	gorm.Model
	SwiftCode     string `gorm:"uniqueIndex;not null"`
	BankName      string `gorm:"not null"`
	Address       string
	CountryISO2   string `gorm:"size:2;not null"`
	CountryName   string `gorm:"not null"`
	IsHeadquarter bool   `gorm:"not null"`
	HeadquarterID *uint  // ID siedziby głównej, jeśli to oddział
}
