# Terminal Path Service
We have used a semi-hexagonal architecture in this microservice. According to the **[[ERD](https://github.com/HeisenGo/heisen-baba/blob/feature/terminal-path/terminal_path/docs/ERD_terminal_path.png.png)]**, there are two entities: Path and Terminal. Based on the architecture, entities are implemented in the storage layer, and their domain-related modules are implemented in the internal layer. Additionally, there are functions that map entity and storage layers.

## Database Models
### Terminal
```go
type Terminal struct {
    gorm.Model
    Name           string `gorm:"type:varchar(100);not null"`
    NormalizedName string `gorm:"type:varchar(100);not null;index:idx_normalized_name_type_city_country,unique,priority:1"`
    Type           string `gorm:"type:varchar(20);not null;index:idx_normalized_name_type_city_country,unique,priority:2"`
    City           string `gorm:"type:varchar(100);not null;index:idx_normalized_name_type_city_country,unique,priority:3"`
    Country        string `gorm:"type:varchar(100);not null;index:idx_normalized_name_type_city_country,unique,priority:4"`
}

```

### Path
```go
type Path struct {
    gorm.Model
    FromTerminalID uint     `gorm:"not null"`
    ToTerminalID   uint     `gorm:"not null"`
    FromTerminal   Terminal `gorm:"foreignKey:FromTerminalID"`
    ToTerminal     Terminal `gorm:"foreignKey:ToTerminalID"`
    DistanceKM     float64  `gorm:"type:decimal(10,2);not null"` // in kilometers
    Code           string   `gorm:"type:varchar(50);not null;uniqueIndex"`
    Name           string   `gorm:"type:varchar(100)"`
    Type           string   `gorm:"type:varchar(20);not null"`
}

```

## Implemented Services
There are several services called by handlers according to our need for CRUD implementation:

### Terminal Services
**Terminal Type Constraints**: Terminal type can only be rail, road, air, or sailing.

### Endpoints
**POST / CreateTerminal**

Creates a terminal based on the TerminalRequest struct in the presenter to handle the POST request for creating a terminal. Validations such as the terminal name are checked in the ops layer.

**GET / GetTerminals**

Fetches all possible terminals according to the field query parameters. The country parameter is mandatory, but city and type parameters are optional.

**PATCH / PatchTerminal**

Updates a terminal based on the terminalID path parameter. If the terminal has related paths, its type, city, and country cannot be updated, but the name can be updated at any time. This logic is implemented in the ops layer.

**DELETE / DeleteTerminal**

Deletes a terminal based on the terminalID path parameter. If there are any paths related to this terminal, it cannot be deleted.

### Path Services
### Endpoints
**POST / CreatePath**

Creates a path based on the PathRequest struct in the presenter to handle the POST request for creating a path. Validations, such as ensuring the same terminal types at the beginning and end of a path, are checked in the ops layer.

**GET / GetPathsByOriginDestinationType**

Fetches all possible paths according to the field query parameters. The from, to, and type parameters, which represent the origin, destination, and type of a path, are optional.

**GET / GetFullPathByID**

Fetches the full path details based on the path ID. This is implemented to be used in other services.

**PATCH / PatchPath**

Updates a path based on the pathID path parameter. If the path has existing unfinished trips, its fromTerminal and toTerminal (and thus its type) cannot be updated, but the code, name, and distance can be updated at any time. This logic is implemented in the ops layer.

**DELETE / DeletePath**

Deletes a path based on the pathID path parameter. If there are any existing unfinished trips, it cannot be deleted.






