# Hotel Room Service

This microservice manages the creation, retrieval, updating, and deletion of hotels and rooms. It includes functionalities for querying and maintaining unique constraints on room names within hotels.

## Database Models

### Hotel
```go
type Hotel struct {
    gorm.Model
    OwnerID   uuid.UUID `gorm:"type:uuid;not null;index:idx_owner_name"`
    Name      string    `gorm:"type:varchar(100);not null;index:idx_owner_name"`
    City      string    `gorm:"type:varchar(100);not null"`
    Country   string    `gorm:"type:varchar(100);not null"`
    Details   string    `gorm:"type:text"`
    IsBlocked bool      `gorm:"default:false"`
    Rooms     []Room
}
Room
go

type Room struct {
    gorm.Model
    HotelID     uint     `gorm:"not null"`
    Hotel       *Hotel   `gorm:"foreignKey:HotelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    Name        string   `gorm:"type:varchar(100);not null;uniqueIndex:idx_hotel_room_name"`
    AgencyPrice uint64   `gorm:"not null"`
    UserPrice   uint64   `gorm:"not null"`
    Facilities  string   `gorm:"type:text"`
    Capacity    uint8    `gorm:"not null"`
    IsAvailable bool     `gorm:"default:true"`
}
Implemented Services
Hotel Services
Hotel Name Constraint: Hotel names are unique per owner, enforced by the unique index on OwnerID and Name.

Endpoints
POST / CreateHotel

Creates a hotel based on the CreateHotelRequest struct. The request includes the hotel details such as name, city, country, and additional information.

GET / GetHotels

Fetches hotels based on optional query parameters such as city, country, and capacity. Pagination is supported through page and page_size parameters.

PUT / UpdateHotel

Updates a hotel based on the id path parameter. The update can include changes to the name, city, country, and blocked status.

DELETE / DeleteHotel

Deletes a hotel based on the id path parameter. Deletion is restricted if the hotel has associated rooms.

Room Services
Unique Room Name Constraint: Room names are unique within each hotel, enforced by a composite unique index.

Endpoints
POST / CreateRoom

Creates a room within a specified hotel based on the CreateRoomRequest struct. Validations include ensuring the uniqueness of the room name within the hotel.

GET / GetRooms

Fetches rooms based on optional query parameters and pagination. Supports filtering by hotel_id and pagination through page and page_size.

PUT / UpdateRoom

Updates a room based on the id path parameter. Optional fields include room name, price, facilities, and availability status.

DELETE / DeleteRoom

Deletes a room based on the id path parameter. Room deletion is straightforward but may be restricted by application logic to maintain consistency.

Validation for Unique Constraints
For maintaining the uniqueness of room names within a hotel, we have applied a composite unique index on the combination of HotelID and Name. This ensures that each room name is unique within the scope of a single hotel.

Validation Logic
Room Creation: When creating a new room, the system checks if the room name already exists for the given hotel. If it does, an error is returned.
Room Update: During room updates, the uniqueness constraint is revalidated to ensure no other room within the same hotel has the same name.
This approach ensures data integrity and prevents duplicate room names within the same hotel, while allowing the same room name to be used in different hotels.