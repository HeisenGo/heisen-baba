url: `localhost:2244/api/v1/terminals`

method: POST

json body sample:


```go
{
  "name": "West Terminal",
  
   "type": "road",
 
   "city": "Tehran",
 
   "country": "Iran"
 
}
```
**GET**

```
localhost:2244/api/v1/terminals?country=Iran

localhost:2244/api/v1/terminals?country=Iran&city=Qazvin

localhost:2244/api/v1/terminals?country=Iran&city=Qazvin&type=air
```

**PATH**
```
localhost:2244/api/v1/paths
```

```go
{

  "from_terminal_id":1,

	"to_terminal_id":6,

"distance_km":120, 

"code":"ZT1R",

"name":"zanjan tehran road"
}

```