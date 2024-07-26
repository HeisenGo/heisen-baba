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


## Terminal City Country Validation
When creating a new terminal or update it, it is indeed a good practice to validate the city and country to ensure data integrity and consistency. We have a couple of options for how to implement this validation using the data from a CSV file whih was provided by **[[simplemaps](https://simplemaps.com/data/world-cities)]**.We filtered the data to only have ascii_city and country. Here are the two main approaches we could use, along with their pros and cons:

### Option 1: Using a Separate SQLite Database
**Advantages:**

- Separation of Concerns: Keeps the city and country data separate from your main application database, which can be cleaner and more modular.
- Independent Updates: Allows you to update the city and country data independently without affecting the main application database.
- Reduced Risk: Any errors or issues in the city/country database do not directly impact the primary database operations.

**Disadvantages:**

- Additional Complexity: Requires maintaining and managing an additional database.
- Slight Performance Overhead: Additional database queries to validate city and country, which could slightly impact performance.

### Option 2: Integrating into the Main Database (Terminal Path DB)
**Advantages:**

- Single Source of Truth: All related data is in one place, simplifying database management and queries.
- Simpler Application Logic: No need to manage connections to multiple databases.
- Potential Performance Gains: Validation queries are within the same database, potentially reducing latency.
  
**Disadvantages:**

- Coupling of Data: City and country data changes can potentially affect the main application database.
- Increased Database Size: The main database grows larger with the additional city and country data.
- Complex Migrations: Schema changes for city/country data may require careful handling to avoid disrupting the main database.

We decided to use **option 1**. cleaning and importing data to sqlight_db is performed using a separate code by python then we integrrated it into config and used the sqlite_db in our app in validation process in internal/terminal/ops.go




