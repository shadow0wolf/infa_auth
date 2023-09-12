package infa_auth

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main1() {
	app := fiber.New()
	app.Use(cors.New())
	//log.Default("starting server")
	api := app.Group("/session-service/api/v1/session")

	// Test handler
	api.Get("/Agent", func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		//log.Fatal(headers)
		h := headers["Ids-Agent-Session-Id"]
		/*
			if h == null {
				log.Println("Ids-Agent-Session-Id header is null")
			}
			h := headers["IDS-AGENT-SESSION-ID"]
			if h == nil {
				log.Println("IDS-AGENT-SESSION-ID header is null")
			}
			h := headers["ids-agent-session-id"]
			if h == nil {
				log.Println("ids-agent-session-id header is null")
			}
		*/
		log.Println("value of header = " + h)

		if "123123123" == h {
			log.Println("SUCCESS")
			return c.SendStatus(200)
		} else {
			log.Println("FAIL")
			return c.SendStatus(401)
		}
	})
	app.Listen(":9898")
}
