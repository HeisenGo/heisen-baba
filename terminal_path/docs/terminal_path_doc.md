


# Terminal Path Service

We have used a semi-hexagonal architecture in this microservice. According to This **[[ERD](https://github.com/HeisenGo/heisen-baba/blob/feature/terminal-path/terminal_path/docs/ERD_terminal_path.png.png)]**. There are two entities ``Path``` and ```Terminal```.
Based on the architecture, entities are implemented in storage layer and their domain-related modules are implemented in internal layer and alse there are some functions that map entity and storage layer.

## Implemented Services:
There are several services called by handlers according to our need for ```CRUD``` implementation:

**Terminal Services**

terminal type can only be rail, road, air, sailing

**POST/GET/PATCH/DELETE**
```
CreateTerminal
```
This function creates a path according to the ```TerminalRequest``` struct which is in the presenter to present post request of creating a terminal. 
The validations like name of the terminal are checked in ops layer.
```
GetTerminals
```
This function will fetch all possible terminals according to the field query parameters, country parameter is nessessary but city and type parameters are optional query parameters.
```
PatchTerminal
```
This request needs a path parameter as terminalID. This is implemented according to this logic if  the terminal has paths related to it its type, city, country cannot be updated but name can be updated in any time. this logic is implemented in ops layer.

```
DeleteTerminal
```
like Patch request takes a terminalID as path parameter Then checks if there is any path related to this terminal if so it can not be deleted




**Path Services**
**POST/GET/PATCH/DELETE**
```
CreatePath
```
This function creates a path according to the ```PathRequest``` struct which is in the presenter to present post request of creating a pth. 
The validations like same terminal types in the begining and end of a path are checked in ops layer.
```
GetPathsByOriginDestinationType
```
This function will fetch all possible paths according to the field query parameters, from and to and type parameters whisch represent the origin, destiantion and type of a path are optional query parameters.
```
GetFullPathByID
```
This is implemented because of the need to be used in other services
```
PatchPath
```
This request a path parameter as pathID. This is implemented according to this logic if  the path has existing unfinished trips its fromTerminal and ToTerminal so its type cannot be updated but code, name, distance can be updated in any time. this logic is implemented in ops layer.
```
DeletePath
```
like Patch request takes a pathID as path parameter Then checks if there is any exstig unfinished trip if so it can not be deleted





