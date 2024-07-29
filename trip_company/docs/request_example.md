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