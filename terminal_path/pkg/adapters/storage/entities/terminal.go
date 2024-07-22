package entities

import (
	"strings"
	"unicode"

	"gorm.io/gorm"
)

type Terminal struct {
	gorm.Model
	Name           string `gorm:"type:varchar(100);not null"`
	NormalizedName string `gorm:"type:varchar(100);not null;index:idx_normalized_name_type_city_country,priority:1"`
	Type           string `gorm:"type:varchar(20);not null;index:idx_name_type_city_country,priority:2"`
	City           string `gorm:"type:varchar(100);not null;index:idx_name_type_city_country,priority:3"`
	Country        string `gorm:"type:varchar(100);not null;index:idx_name_type_city_country,priority:4"`
	OutgoingPaths  []Path `gorm:"foreignKey:FromTerminalID"`
	IncomingPaths  []Path `gorm:"foreignKey:ToTerminalID"`

	// Composite unique constraint
	UniqueTerminal string `gorm:"uniqueIndex:idx_name_type_city_country"`
}

func NormalizeName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Remove extra spaces
	name = strings.Join(strings.Fields(name), " ")

	// Remove all non-alphanumeric characters except spaces
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == ' ' {
			return r
		}
		return -1
	}, name)
}

func (t *Terminal) BeforeCreate(tx *gorm.DB) (err error) {
	t.NormalizedName = NormalizeName(t.Name)
	return nil
}

func (t *Terminal) BeforeUpdate(tx *gorm.DB) (err error) {
	t.NormalizedName = NormalizeName(t.Name)
	return nil
}
