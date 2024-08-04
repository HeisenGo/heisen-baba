# heisen-baba / Trip Company Microservice

If you need to independently run this service. You need to use ```independentServiceconfig.yaml.example``` rename it to config.yaml then comment bellow command in service/app.go.
```
// service registry
	app.mustRegisterService(cfg.Server)
```

after setting up your database run:
```
go run cmd/main.go --config config.yaml
```

for more info about implementation check docs folder
