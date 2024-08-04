# Hotel and Room API Documentation

## Hotels

### Create Hotel
**URL:** `localhost:8080/api/v1/hotels`  
**Method:** POST  

**Request Body:**
```json
{
  "name": "Grand Hotel",
  "city": "New York",
  "country": "USA",
  "details": "A luxurious hotel in the heart of the city.",
  "is_blocked": false
}
```
### Get Hotels
**URL:** `localhost:8080/api/v1/hotels`
**Method:** GET

**Query Parameters:**

**city -**   Filter hotels by city.
**country -**   Filter hotels by country.
**capacity -**   Filter hotels by room capacity.
**page -**   Pagination: page number.
**page_size -**   Pagination: number of items per page.

**Example Requests:**
`localhost:8080/api/v1/hotels?country=USA`
`localhost:8080/api/v1/hotels?country=USA&city=New%20York`
`localhost:8080/api/v1/hotels?country=USA&city=New%20York&capacity=5`
### Get Hotels by Owner ID
**URL:** localhost:8080/api/v1/hotels/owner
**Method:** GET

### Query Parameters:
**owner_id -** Filter hotels by owner ID.
**page -** Pagination: page number.
**page_size -** Pagination: number of items per page.

**Example Requests:**
`localhost:8080/api/v1/hotels/owner?owner_id=123e4567-e89b-12d3-a456-426614174000`
`localhost:8080/api/v1/hotels/owner?owner_id=123e4567-e89b-12d3-a456-426614174000&page=1&page_size=10`
### Update Hotel
**URL:** `localhost:8080/api/v1/hotels/{id}`
**Method:** PUT

**Request Body:**

```json
{
  "name": "Updated Grand Hotel",  // required
  "city": "Los Angeles",          // required
  "country": "USA",               // required
  "details": "Updated details.",  // required
  "is_blocked": true              // required
}
```
**Example Request:**
`localhost:8080/api/v1/hotels/5`

### Delete Hotel
**URL**: `localhost:8080/api/v1/hotels/{id}`
**Method**: DELETE

**Example Request:**
`localhost:8080/api/v1/hotels/3`


## Rooms
### Create Room
**URL:** `localhost:8080/api/v1/rooms`
**Method:**  POST

**Request Body:**
```json
{
  "hotel_id": 1,
  "name": "Deluxe Suite",
  "agency_price": 15000,
  "user_price": 20000,
  "facilities": "WiFi, TV, Minibar",
  "capacity": 2,
  "is_available": true
}
```
### Get Rooms
**URL:**  `localhost:8080/api/v1/rooms`
**Method:** GET

**Query Parameters:**
**page -**  Pagination: page number.
**page_size -** Pagination: number of items per page.

**Example Requests:**
`localhost:8080/api/v1/rooms`
`localhost:8080/api/v1/rooms?page=1&page_size=10`

### Update Room
**URL:** ` localhost:8080/api/v1/rooms/{id}`
**Method:** PUT

**Request Body:**

```json
{
  "name": "Updated Suite",        // optional
  "agency_price": 16000,          // optional
  "user_price": 21000,            // optional
  "facilities": "WiFi, TV, Minibar, Pool", // optional
  "capacity": 3,                  // optional
  "is_available": false           // optional
}
```
**Example Request:**
`localhost:8080/api/v1/rooms/10`

### Delete Room
**URL:**  `localhost:8080/api/v1/rooms/{id}`
**Method:** DELETE

**Example Request:**
`localhost:8080/api/v1/rooms/7`