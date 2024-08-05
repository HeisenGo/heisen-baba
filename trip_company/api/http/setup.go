package http

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"tripcompanyservice/api/http/handlers"
	middlerwares "tripcompanyservice/api/http/middlewares"
	"tripcompanyservice/config"
	adapter "tripcompanyservice/pkg/adapters"
	"tripcompanyservice/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Run(cfg config.Server, app *service.AppContainer) {
	fiberApp := fiber.New()
	api := fiberApp.Group("/api/v1") //, middlerwares.SetUserContext())

	createGroupLogger := loggerSetup(fiberApp)

	registerGlobalRoutes(api)
	registerTransportCompanyRoutes(api, app, createGroupLogger("companies"))

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.HttpPort)))
}

func loggerSetup(app *fiber.App) func(groupName string) fiber.Handler {

	// Create the logs directory if it does not exist
	logDir := "./logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("error creating logs directory: %v", err)
	}

	// Common format for console logging
	consoleLoggerConfig := logger.Config{
		Format:     "${time} [${ip}]:${port} ${status} - ${method} ${path} - ${latency} ${bytesSent} ${bytesReceived} ${userAgent}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Local",
	}
	app.Use(logger.New(consoleLoggerConfig))

	// Function to create a logger middleware with separate log file
	createGroupLogger := func(groupName string) fiber.Handler {
		logFilePath := filepath.Join(logDir, groupName+".log")
		file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}

		return logger.New(logger.Config{
			Format:     "${time} [${ip}]:${port} ${status} - ${method} ${path} - ${latency} ${bytesSent} ${bytesReceived} ${userAgent}\n",
			TimeFormat: "02-Jan-2006 15:04:05",
			TimeZone:   "Local",
			Output:     file,
		})
	}
	return createGroupLogger
}

func registerGlobalRoutes(router fiber.Router) {
	// Setup a simple health check route
	router.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Service is up and running")
	})
}

func registerTransportCompanyRoutes(router fiber.Router, app *service.AppContainer, loggerMiddleWare fiber.Handler) {
	router = router.Group("/companies")
	router.Use(loggerMiddleWare)

	router.Post("",
		//middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		//middlewares.Auth(),
		handlers.CreateTransportCompany(app.CompanyService()),
	)
	router.Get("/my-companies/:ownerID",
		//middlewares.Auth(),
		handlers.GetUserCompanies(app.CompanyService()),
	)
	router.Get("",
		//middlewares.Auth(),
		handlers.GetCompanies(app.CompanyService()),
	)

	// only owner can do this
	router.Delete("/my-companies/:companyID",
		handlers.DeleteCompany(app.CompanyService()),
	)
	router.Patch("/my-companies/:companyID",
		handlers.PatchCompany(app.CompanyService()),
	)

	// only admin can do this
	router.Patch("/block/:companyID", //, middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
		//middlewares.Auth(),
		handlers.BlockCompany(app.CompanyService()))

	router.Post("/trips",
		//middlewares.Auth(),
		middlerwares.SetTransaction(adapter.NewGormCommiter(app.RawDBConnection())),
		//handlers.CreateTrip(app.TripServiceFromCtx)),
		handlers.CreateTrip(app.TripServiceFromCtx))

	router.Get("/trips", handlers.GetTrips(app.TripService()))
	router.Get("/one-trip/:tripID", handlers.GetFullTripByID(app.TripService()))
	router.Patch("/trips/:tripID",
		handlers.PatchTrip(app.TripService()),
	)
	router.Get("/company-trips/:companyID", handlers.GetCompanyTrips(app.TripService()))
	router.Post("/buy", handlers.BuyTicket(app.TicketService())) // TODO : should be transactional

	router.Patch("/cancel-ticket/:ticketID", handlers.CancelTicketByID(app.TicketService()))
	router.Get("/tickets", handlers.GetTickets(app.TicketService())) // configure a way o separate agnecy and user TODO:

	router.Post("/vehicle-req", handlers.CreateVehicleRequest(app.VehicleReqService()))

	router.Post("/tech-teams", handlers.CreateTechTeam(app.TechTeamService()))
	router.Post("/tech-members", handlers.CreateTechMember(app.TechTeamService()))
	router.Get("/tech-teams/:companyID", handlers.GetTechTeamsOfCompany(app.TechTeamService()))
	router.Patch("/set-team/:tripID", handlers.SetTechTeamToTrip(app.TripServiceFromCtx))
	router.Patch("/cancel-trip/:tripID", handlers.CancelTrip(app.TripServiceFromCtx))
	router.Patch("/finish-trip/:tripID", handlers.FinishTrip(app.TripServiceFromCtx))
	router.Patch("/confirm-trip/:tripID", handlers.ConfirmTrip(app.TripServiceFromCtx))

	router.Get("/path-trips/:pathID", handlers.GetCountPathUnfinishedTrips(app.TripService()))
}
