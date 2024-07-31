# Request Samples

# Transport Company

## POST: Create
url: ```localhost:2245/api/v1/companies```
```go
{
"name": "mahan 5" ,
"desc": "we support air and road tickets around the world", //optional
"owner_id": 5 ,
	"address": "Tehran Azadi Street" //optional
}

```
## GET companies(admin) or my companies(admin/owner)

```
localhost:2245/api/v1/companies
```
```
localhost:2245/api/v1/companies/my-companies/owner_ID
```

## DELETE : only by owner
```
localhost:2245/api/v1/companies/my-companies/3
```

## Block UnBlock only by admin
```
localhost:2245/api/v1/companies/block/3
```

```go
{
  "is_blocked":true
}
```

## PATH only owner
```
localhost:2245/api/v1/companies/my-companies/3
```

```go
{
  "name": "salman",
  "desc": "gh ",
  "address": "new"
  //"new_owner_email": 4 
}

```

# Trip

## CREATE
``` localhost:2245/api/v1/companies/trips```
```go
{
  "company_id": 1,
  "user_date": "2024-08-30 23:27:09",
  "tour_date": "2024-08-25 23:27:09",
  "user_price": 45000,
  "agency_price": 30000,
  "path_id": 2,
  "min_pass": 20,
  //"tech_team_id": 1,
  "max_tickets": 23,
  "start_date": "2024-09-30 23:27:09",
  "penalty": {
    "first_days": 25,
    "first_percentage": 10,
    "second_days": 15,
    "second_percentage": 35,
    "third_days": 10,
    "third_percentage": 45
  }
}
```