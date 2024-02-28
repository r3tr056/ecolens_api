package models

import (
	"gorm.io/gorm"
)

type Brand struct {
	gorm.Model
	Name string `json:"name"`
}

type ProductImage struct {
	gorm.Model
	ProductID int    `json:"product_id"`
	Image     string `json:"image"`
}

type Category struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Products    []Product `gorm:"foreignKey:CategoryID"`
}

type LCAMetrics struct {
	gorm.Model
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
	EPDID uint    `json:"epd_id"`
}

type EnvironmentalProductDeclaration struct {
	gorm.Model
	ProductID   uint         `json:"product_id"`
	Description string       `json:"description"`
	LCAMetrics  []LCAMetrics `gorm:"foreignKey:EPDID" json:"lca_metrics"`
}

type Report struct {
	gorm.Model
	Name    string `json:"name"`
	EPDID   uint   `json:"epd_id"`
	Summary string `json:"summary"`
}

type Product struct {
	gorm.Model
	Name            string                          `gorm:"type:varchar(255);index:idx_name_gin"`
	BrandID         uint                            `json:"brand_id"`
	Barcode         string                          `json:"barcode"`
	Brand           Brand                           `json:"brand"`
	Images          []ProductImage                  `json:"marketplace_alternatives"`
	EPD             EnvironmentalProductDeclaration `json:"epd"`
	Reports         []Report                        `gorm:"foreignKey:EPDID"`
	Description     string                          `gorm:"type:text;index:idx_description_gin"`
	Price           float64                         `json:"price"`
	Link            string                          `json:"link"`
	CategoryID      uint
	Category        Category
	EnvironmentTags []EnvironmentTag `gorm:"foreignKey:ProductID"`
}

type MarketPlaceProduct struct {
	gorm.Model
	Product
	Stock             int    `json:"stock"`
	Price             string `json:"price"`
	Desc              string `json:"desc"`
	PubDate           string `json:"pub_date"`
	BrandID           uint   `json:"brand_id"`
	Quantity          int    `json:"quantity"`
	BarCode           string `json:"bar_code"`
	ExpiryDate        string `json:"expiry_date"`
	DateOfManufacture string `json:"date_of_manufacture"`
}

type EnvironmentTag struct {
	gorm.Model
	Name      string `json:"name"`
	ProductID uint
	Product   Product `gorm:"foreignKey:ProductID"`
}
