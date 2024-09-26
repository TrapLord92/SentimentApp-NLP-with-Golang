package main

import (
	"fmt"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"github.com/grassmudhorses/vader-go/lexicon"
	"github.com/grassmudhorses/vader-go/sentitext"
)

// SentimentDetails struct holds the sentiment analysis result fields.
type SentimentDetails struct {
	Positive float64 `json:"positive"`
	Negative float64 `json:"negative"`
	Neutral  float64 `json:"neutral"`
	Compound float64 `json:"compound"`
}

func main() {
	// Initialize HTML template engine for views
	engine := html.New("./views", ".html") // Fiber's latest HTML template engine
	engine.Reload(true)                    // Enable hot-reload of templates for development

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Route: Render the home page
	app.Get("/", func(c *fiber.Ctx) error {
		initMessage := "Sentiment Analysis App on the Go"
		return c.Render("index", fiber.Map{
			"initMessage": initMessage,
		})
	})

	// Route: Handle POST form data for sentiment analysis
	app.Post("/", func(c *fiber.Ctx) error {
		initMessage := "Sentiment Analysis App on the Go"
		message := c.FormValue("message")

		// Handle case where no message is provided
		if message == "" {
			return c.Render("index", fiber.Map{
				"initMessage": "Please provide a message for analysis.",
			})
		}

		sentimentResults := sentimentize(message)
		return c.Render("index", fiber.Map{
			"initMessage":      initMessage,
			"originalMsg":      message,
			"sentimentDetails": sentimentResults,
		})
	})

	// API Route: Analyze sentiment of the given text in the query parameter
	app.Get("/api", func(c *fiber.Ctx) error {
		message := c.Query("text")

		// If no text is provided, respond with an error
		if message == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No text provided for sentiment analysis",
			})
		}

		sentimentResults := analyzeSentiment(message)
		return c.JSON(fiber.Map{
			"message":   message,
			"sentiment": sentimentResults,
		})
	})

	// Swagger Route: Serve Swagger documentation
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "http://example.com/doc.json", // Custom URL for Swagger JSON
		DeepLinking: false,
	}))

	// Start the server and listen on port 3000
	if err := app.Listen(":3000"); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

// Sentiment analysis for the web interface, returning detailed scores
func sentimentize(docx string) SentimentDetails {
	// Parse the text using VADER's lexicon
	parsedtext := sentitext.Parse(docx, lexicon.DefaultLexicon)
	// Calculate sentiment polarity scores
	results := sentitext.PolarityScore(parsedtext)

	// Return sentiment details
	return SentimentDetails{
		Positive: results.Positive,
		Negative: results.Negative,
		Neutral:  results.Neutral,
		Compound: results.Compound,
	}
}

// Analyze sentiment for the API, returning only the compound score
func analyzeSentiment(docx string) float64 {
	// Parse the text using VADER's lexicon
	parsedtext := sentitext.Parse(docx, lexicon.DefaultLexicon)
	// Calculate sentiment polarity scores
	results := sentitext.PolarityScore(parsedtext)

	// Return only the compound score
	return results.Compound
}
