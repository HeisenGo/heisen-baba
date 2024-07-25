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

**Terminal Patch**
```localhost:2244/api/v1/terminals/5```

```go
{
  "name": "tt", //optional
	"type": "rail", // optional
	"city": "Hamedan", // optional
	"country": "Iran" //optional
}

```
**Terminal Delete**

``` localhost:2244/api/v1/terminals/3 ```


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

**GET**
```
localhost:2244/api/v1/paths

localhost:2244/api/v1/paths?from=Tehran

localhost:2244/api/v1/paths?from=Tehran&to=Zanjan

localhost:2244/api/v1/paths?from=Tehran&to=Zanjan&type=road
```

**Patch Path**
```localhost:2244/api/v1/paths/3```

```go
{
  "from_terminal_id":14, //optional

 "to_terminal_id":18, //optional

"distance_km":810,  //optional

"code":"TMR1", //optional

"name":"updated tehran mashhad road 1" //optional
}
```

**delete Path**
```
localhost:2244/api/v1/paths/4
```