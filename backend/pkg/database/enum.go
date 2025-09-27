package database

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

type Enum struct {
	Name   string
	Values []string
}

func SeedEnums(db *gorm.DB) {
	enums := []Enum{}

	for _, enum := range enums {
		dropQuery := fmt.Sprintf(`DROP TYPE IF EXISTS "%s" CASCADE;`, enum.Name)
		if err := db.Exec(dropQuery).Error; err != nil {
			log.Fatalf("Error dropping enum %s: %v", enum.Name, err)
		}

		// Sonra enum'u yeniden oluştur
		createQuery := fmt.Sprintf(`
	DO $$
	BEGIN
		CREATE TYPE "%s" AS ENUM (%s);
	END $$;`,
			enum.Name,
			FormatEnumValues(enum.Values),
		)

		if err := db.Exec(createQuery).Error; err != nil {
			log.Fatalf("Error creating enum %s: %v", enum.Name, err)
		}

		query := fmt.Sprintf(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 
				FROM pg_type t
				WHERE t.typname = '%s'
			) THEN
				CREATE TYPE "%s" AS ENUM (%s);
			END IF;
		END $$;`,
			enum.Name,
			enum.Name,
			FormatEnumValues(enum.Values),
		)

		if err := db.Exec(query).Error; err != nil {
			log.Fatalf("Error creating enum %s: %v", enum.Name, err)
		}
	}
}

func FormatEnumValues(values []string) string {
	formattedValues := make([]string, len(values))
	for i, v := range values {
		formattedValues[i] = fmt.Sprintf("'%s'", v)
	}
	return strings.Join(formattedValues, ", ")
}
