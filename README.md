# rest
    import "github.com/tmc/rest"

Package rest provides a micro framework for writing json HTTP endpoints

Inspiration from dougblack's sleepy

Example:


```go
	package main
	
	import (
		"fmt"
		"io/ioutil"
		"net/http"
		"net/url"
	
		"github.com/tmc/rest"
	)
	
	type User struct {
		Name string
	}
	
	type UserRepository struct {
		users []User
	}
	
	func (ur UserRepository) AllUsers() []User {
		return ur.users
	}
	
	type UserList struct {
		repo UserRepository
	}
	
	func (ul UserList) Get(values url.Values) (int, interface{}) {
		return http.StatusOK, ul.repo.AllUsers()
	}
	
	func main() {
		a := rest.API{}
		a.AddResource("/users", UserList{UserRepository{[]User{{"joe"}, {"sally"}}}})
	
		go a.Start(8080)
	
		resp, err := http.Get("http://127.0.0.1:8080/users")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
	
		fmt.Println(resp.Status, string(body))
        // Output:
    	// 200 OK [{"Name":"joe"},{"Name":"sally"}]
	}
```

## Constants
``` go
const (
    GET    = "GET"
    POST   = "POST"
    PUT    = "PUT"
    DELETE = "DELETE"
)
```

## type API
``` go
type API struct{}
```

### func (\*API) Abort
``` go
func (api *API) Abort(rw http.ResponseWriter, statusCode int)
```


### func (\*API) AddResource
``` go
func (api *API) AddResource(path string, resource Resource)
```


### func (\*API) Start
``` go
func (api *API) Start(port int)
```


## type BaseResource
``` go
type BaseResource struct{}
```


### func (BaseResource) Delete
``` go
func (BaseResource) Delete(values url.Values) (int, interface{})
```


### func (BaseResource) Get
``` go
func (BaseResource) Get(values url.Values) (int, interface{})
```


### func (BaseResource) Post
``` go
func (BaseResource) Post(values url.Values) (int, interface{})
```


### func (BaseResource) Put
``` go
func (BaseResource) Put(values url.Values) (int, interface{})
```


## type Resource
``` go
type Resource interface{}
```

