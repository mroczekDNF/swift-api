package models

// SwiftCode modeluje dane dla kodów SWIFT
// Usunięto GORM i dostosowano do database/sql

type SwiftCode struct {
	ID            int64  // Unikalny identyfikator (zamiast gorm.Model)
	SwiftCode     string // Kod SWIFT (unikalny, niepusty)
	BankName      string // Nazwa banku (niepusty)
	Address       string // Adres banku
	CountryISO2   string // Kod kraju (ISO2, zawsze 2 znaki)
	CountryName   string // Nazwa kraju (niepusty)
	IsHeadquarter bool   // Czy to siedziba główna
	HeadquarterID *int64 // ID siedziby głównej (dla oddziałów)
}
