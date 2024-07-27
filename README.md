# heisen-baba / Terminal Path Microservice

If you need to independently run this service. You need to use ```independentServiceconfig.yaml.example``` rename it to config.yaml then comment bellow command in service/app.go.
```
// service registry
	app.mustRegisterService(cfg.Server)
```

after setting up your database run:
```
go run cmd/main.go --config config.yaml
```

for more info about implementation check docs/terminal_path_doc.md.

### Terminal_Path **[[link](https://github.com/HeisenGo/heisen-baba/blob/feature/terminal-path/terminal_path/docs/terminal_path_doc.md)]**
