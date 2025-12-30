package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/entorno35/backend/internal/config"
	"github.com/entorno35/backend/internal/database"
	"github.com/entorno35/backend/internal/domain"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Get database connection string
	dsn, err := cfg.DatabaseURL()
	if err != nil {
		log.Fatalf("Failed to get database URL: %v", err)
	}

	// Connect to database
	db, err := database.Connect(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("✅ Database connected successfully")

	// Read JSON file
	jsonFile := "nom035_questions.json"
	if len(os.Args) > 1 {
		jsonFile = os.Args[1]
	}

	data, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read JSON file %s: %v", jsonFile, err)
	}

	var questionsData domain.QuestionsDataJSON
	if err := json.Unmarshal(data, &questionsData); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	log.Printf("✅ Parsed JSON file: %s", jsonFile)

	// Seed each guide
	if err := seedGuide(db, questionsData.GuideI, domain.GuideTypeI); err != nil {
		log.Fatalf("Failed to seed Guide I: %v", err)
	}
	log.Println("✅ Guide I seeded successfully")

	if err := seedGuide(db, questionsData.GuideII, domain.GuideTypeII); err != nil {
		log.Fatalf("Failed to seed Guide II: %v", err)
	}
	log.Println("✅ Guide II seeded successfully")

	if err := seedGuide(db, questionsData.GuideIII, domain.GuideTypeIII); err != nil {
		log.Fatalf("Failed to seed Guide III: %v", err)
	}
	log.Println("✅ Guide III seeded successfully")

	log.Println("🎉 Database seeding completed successfully!")
}

// seedGuide seeds a single guide into the database
func seedGuide(db *gorm.DB, guide domain.GuideJSON, guideType domain.GuideType) error {
	questionType := domain.QuestionTypeBinary
	if guide.Type == "likert" {
		questionType = domain.QuestionTypeLikert
	}

	for i, qJSON := range guide.Questions {
		var categoryID, domainID, dimensionID *uint
		var polarity *domain.QuestionPolarity
		var section, subsection *string

		if guideType == domain.GuideTypeI {
			// Guide I: Map Section -> Category, Subsection -> Domain (create for consistency but don't link via FK)
			section = &qJSON.Section
			subsection = &qJSON.Subsection
			// Guide I doesn't have polarity
			polarity = nil
			
			// Create Category/Domain entries for Guide I (for organizational purposes)
			// but don't link Question via FK (category_id/domain_id remain NULL)
			var category domain.Category
			if err := db.Where("name = ?", qJSON.Section).FirstOrCreate(&category, domain.Category{
				Name: qJSON.Section,
			}).Error; err != nil {
				return fmt.Errorf("failed to find or create category %s: %w", qJSON.Section, err)
			}

			var domainObj domain.Domain
			if err := db.Where("category_id = ? AND name = ?", category.ID, qJSON.Subsection).
				FirstOrCreate(&domainObj, domain.Domain{
					CategoryID: category.ID,
					Name:       qJSON.Subsection,
				}).Error; err != nil {
				return fmt.Errorf("failed to find or create domain %s: %w", qJSON.Subsection, err)
			}

			// Guide I questions don't link to category/domain via FK (use section/subsection fields instead)
			categoryID = nil
			domainID = nil
			dimensionID = nil
		} else {
			// Guide II/III: Build normalized hierarchy
			section = nil
			subsection = nil
			
			// Parse polarity
			if qJSON.Polarity != "" {
				p := domain.QuestionPolarity(qJSON.Polarity)
				polarity = &p
			}

			// Find or create Category
			var category domain.Category
			if err := db.Where("name = ?", qJSON.Category).FirstOrCreate(&category, domain.Category{
				Name: qJSON.Category,
			}).Error; err != nil {
				return fmt.Errorf("failed to find or create category %s: %w", qJSON.Category, err)
			}
			categoryID = &category.ID

			// Find or create Domain (must have category_id)
			var domainObj domain.Domain
			if err := db.Where("category_id = ? AND name = ?", category.ID, qJSON.Domain).
				FirstOrCreate(&domainObj, domain.Domain{
					CategoryID: category.ID,
					Name:       qJSON.Domain,
				}).Error; err != nil {
				return fmt.Errorf("failed to find or create domain %s: %w", qJSON.Domain, err)
			}
			domainID = &domainObj.ID

			// Find or create Dimension (if provided)
			if qJSON.Dimension != "" {
				var dimension domain.Dimension
				if err := db.Where("domain_id = ? AND name = ?", domainObj.ID, qJSON.Dimension).
					FirstOrCreate(&dimension, domain.Dimension{
						DomainID: domainObj.ID,
						Name:     qJSON.Dimension,
					}).Error; err != nil {
					return fmt.Errorf("failed to find or create dimension %s: %w", qJSON.Dimension, err)
				}
				dimensionID = &dimension.ID
			}
		}

		// Create or update Question
		question := domain.Question{
			QuestionNumber: qJSON.Number,
			GuideType:      guideType,
			Type:           questionType,
			Text:           qJSON.Text,
			Polarity:       polarity,
			Section:        section,
			Subsection:     subsection,
			CategoryID:     categoryID,
			DomainID:       domainID,
			DimensionID:    dimensionID,
			OrderIndex:     i,
		}

		// Check if question already exists (by guide_type and question_number)
		var existingQuestion domain.Question
		result := db.Where("guide_type = ? AND question_number = ?", guideType, qJSON.Number).
			First(&existingQuestion)

		if result.Error == gorm.ErrRecordNotFound {
			// Create new question
			if err := db.Create(&question).Error; err != nil {
				return fmt.Errorf("failed to create question %d: %w", qJSON.Number, err)
			}
			log.Printf("  Created question %d: %s", qJSON.Number, truncateString(qJSON.Text, 60))
		} else if result.Error != nil {
			return fmt.Errorf("failed to check for existing question %d: %w", qJSON.Number, result.Error)
		} else {
			// Update existing question
			question.ID = existingQuestion.ID
			if err := db.Model(&existingQuestion).Updates(question).Error; err != nil {
				return fmt.Errorf("failed to update question %d: %w", qJSON.Number, err)
			}
			log.Printf("  Updated question %d: %s", qJSON.Number, truncateString(qJSON.Text, 60))
		}
	}

	return nil
}

// truncateString truncates a string to the specified length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

