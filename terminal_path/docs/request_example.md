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