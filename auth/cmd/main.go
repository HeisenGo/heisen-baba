package main

import (
	http_server "authservice/api/http"
	"authservice/config"
	"authservice/service"
	"flag"
	"log"
	"os"
)

var configPath = flag.String("config", "", "configuration path")

func main() {
	cfg := readConfig()

	app, err := service.NewAppContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	http_server.Run(cfg, app)
}

func readConfig() config.Config {
	flag.Parse()

	if cfgPathEnv := os.Getenv("APP_CONFIG_PATH"); len(cfgPathEnv) > 0 {
		*configPath = cfgPathEnv
	}

	if len(*configPath) == 0 {
		log.Fatal("configuration file not found")
	}

	cfg, err := config.ReadStandard(*configPath)

	if err != nil {
		log.Fatal(err)
	}

	return cfg
}


package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/your-repo/auth/api/http/handlers"
    "github.com/your-repo/auth/api/http/middlewares"
    "github.com/your-repo/auth/config"
    "github.com/your-repo/auth/internal/user"
    "github.com/your-repo/auth/pkg/adapters/consul"
    "github.com/your-repo/auth/pkg/adapters/storage"
    "github.com/your-repo/auth/pkg/rbac"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Initialize database connection
    db, err := storage.SetupDatabase(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to set up database: %v", err)
    }
    defer db.Close()

    // Initialize storage
    storageAdapter := storage.NewStorage(db)

    // Initialize user service
    userService := user.NewService(storageAdapter)

    // Initialize RBAC
    rbac.InitRBAC()

    // Create router
    router := mux.NewRouter()

    // Create auth handler
    authHandler := handlers.NewAuthHandler(userService)

    // Set up routes
    router.HandleFunc("/register", authHandler.Register).Methods("POST")
    router.HandleFunc("/login", authHandler.Login).Methods("POST")

    // Protected routes
    protectedRouter := router.PathPrefix("/api").Subrouter()
    protectedRouter.Use(middlewares.AuthMiddleware)

    // Example protected routes
    protectedRouter.HandleFunc("/company", middlewares.RBACMiddleware("view_company")(handlers.ViewCompany)).Methods("GET")
    protectedRouter.HandleFunc("/company", middlewares.RBACMiddleware("edit_company")(handlers.EditCompany)).Methods("PUT")
    protectedRouter.HandleFunc("/trip", middlewares.RBACMiddleware("view_trip")(handlers.ViewTrip)).Methods("GET")
    protectedRouter.HandleFunc("/trip", middlewares.RBACMiddleware("edit_trip")(handlers.EditTrip)).Methods("PUT")

    // Health check endpoint for Consul
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "OK")
    }).Methods("GET")

    // Register service with Consul
    serviceID := fmt.Sprintf("%s-%s", cfg.ServiceName, os.Getenv("HOSTNAME"))
    err = consul.RegisterService(cfg.ServiceName, serviceID, cfg.ServiceHost, cfg.ServicePort)
    if err != nil {
        log.Printf("Failed to register service with Consul: %v", err)
    }

    // Start the server
    addr := fmt.Sprintf("%s:%d", cfg.ServiceHost, cfg.ServicePort)
    log.Printf("Starting server on %s", addr)
    log.Fatal(http.ListenAndServe(addr, router))
}