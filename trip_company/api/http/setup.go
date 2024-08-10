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
	api := fiberApp.Group("/api/v1",middlerwares.SetUserContext()) //, middlerwares.SetUserContext())

	createGroupLogger := loggerSetup(fiberApp)

	registerGlobalRoutes(api)
	registerTransportCompanyRoutes(api, app, createGroupLogger("companies"))

	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.HTTPPort)))
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
		middlerwares.Auth(app.AuthClient()),
		handlers.CreateTransportCompany(app.CompanyService()),
	)
	router.Get("/my-companies/:ownerID",
	middlerwares.Auth(app.AuthClient()),
		handlers.GetUserCompanies(app.CompanyService()),
	)
	router.Get("",
	middlerwares.Auth(app.AuthClient()),
			handlers.GetCompanies(app.CompanyService()),
	)

	// only owner can do this
	router.Delete("/my-companies/:companyID",
	middlerwares.Auth(app.AuthClient()),
		handlers.DeleteCompany(app.CompanyService()),
	)
	router.Patch("/my-companies/:companyID",
	middlerwares.Auth(app.AuthClient()),
		handlers.PatchCompany(app.CompanyService()),
	)

	// only admin can do this
	router.Patch("/block/:companyID", //, middlewares.SetTransaction(adapters.NewGormCommitter(app.RawDBConnection())),
	middlerwares.Auth(app.AuthClient()),
		handlers.BlockCompany(app.CompanyService()))

	router.Post("/trips",
		//middlewares.Auth(),
		middlerwares.SetTransaction(adapter.NewGormCommiter(app.RawDBConnection())),
		middlerwares.Auth(app.AuthClient()),
				handlers.CreateTrip(app.TripServiceFromCtx))

	router.Get("/trips", middlerwares.Auth(app.AuthClient()),handlers.GetTrips(app.TripService()))
	router.Get("/agency-trips", middlerwares.Auth(app.AuthClient()),handlers.GetAgencyTrips(app.TripService()))

	router.Get("/one-trip/:tripID", middlerwares.Auth(app.AuthClient()),handlers.GetFullTripByID(app.TripService()))
	router.Get("/one-agency-trip/:tripID",middlerwares.Auth(app.AuthClient()), handlers.GetFullAgencyTripByID(app.TripService()))

	router.Patch("/trips/:tripID",middlerwares.Auth(app.AuthClient()),
		handlers.PatchTrip(app.TripServiceFromCtx),
	)
	router.Get("/company-trips/:companyID", middlerwares.Auth(app.AuthClient()),handlers.GetCompanyTrips(app.TripService()))
	router.Get("/company-agency-trips/:companyID",middlerwares.Auth(app.AuthClient()), handlers.GetCompanyAgencyTrips(app.TripService()))


	router.Post("/buy", handlers.BuyTicket(app.TicketServiceFromCtx))
	router.Patch("/cancel-ticket/:ticketID",middlerwares.Auth(app.AuthClient()), handlers.CancelTicketByID(app.TicketServiceFromCtx))

	router.Get("/user-tickets", middlerwares.Auth(app.AuthClient()),handlers.GetUserTickets(app.TicketService()))
	router.Get("/agency-tickets/:agencyID", middlerwares.Auth(app.AuthClient()),handlers.GetAgencyTickets(app.TicketService())) 

	router.Post("/vehicle-req", middlerwares.Auth(app.AuthClient()),handlers.CreateVehicleRequest(app.VehicleReqServiceFromCtx))
	router.Delete("/vehicle-req:vRID",middlerwares.Auth(app.AuthClient()), handlers.DeleteVR(app.VehicleReqServiceFromCtx))

	router.Post("/tech-teams", middlerwares.Auth(app.AuthClient()),handlers.CreateTechTeam(app.TechTeamService()))
	router.Delete("/tech-teams/:teamID",middlerwares.Auth(app.AuthClient()), handlers.DeleteTeam(app.TechTeamService()))

	router.Post("/tech-members", middlerwares.Auth(app.AuthClient()),handlers.CreateTechMember(app.TechTeamService()))
	router.Get("/tech-teams/:companyID", middlerwares.Auth(app.AuthClient()),handlers.GetTechTeamsOfCompany(app.TechTeamService()))
	router.Patch("/set-team/:tripID", middlerwares.Auth(app.AuthClient()),handlers.SetTechTeamToTrip(app.TripServiceFromCtx))
	router.Patch("/cancel-trip/:tripID",middlerwares.Auth(app.AuthClient()), handlers.CancelTrip(app.TripServiceFromCtx))
	router.Patch("/finish-trip/:tripID",middlerwares.Auth(app.AuthClient()), handlers.FinishTrip(app.TripServiceFromCtx))
	router.Patch("/confirm-trip/:tripID",middlerwares.Auth(app.AuthClient()), handlers.ConfirmTrip(app.TripServiceFromCtx))

	router.Get("/path-trips/:pathID", handlers.GetCountPathUnfinishedTrips(app.TripService()))
}
