// Package rest provides a micro framework for writing json HTTP endpoints
//
// Inspiration from dougblack's sleepy
//
// Example:
//  package main
//  
//  import (
//  	"fmt"
//  	"io/ioutil"
//  	"net/http"
//  	"net/url"
//  
//  	"github.com/tmc/rest"
//  )
//  
//  type User struct {
//  	Name string
//  }
//  
//  type UserRepository struct {
//  	users []User
//  }
//  
//  func (ur UserRepository) AllUsers() []User {
//  	return ur.users
//  }
//  
//  type UserList struct {
//  	repo UserRepository
//  }
//  
//  func (ul UserList) Get(values url.Values) (int, interface{}) {
//  	return http.StatusOK, ul.repo.AllUsers()
//  }
//  
//  func main() {
//  	a := rest.API{}
//  	a.AddResource("/users", UserList{UserRepository{[]User{{"joe"}, {"sally"}}}})
//  
//  	go a.Start(8080)
//  
//  	resp, err := http.Get("http://127.0.0.1:8080/users")
//  	if err != nil {
//  		panic(err)
//  	}
//  	defer resp.Body.Close()
//  	body, err := ioutil.ReadAll(resp.Body)
//  
//  	fmt.Println(resp.Status, string(body))
//  	// Output:
//  	// 200 OK [{"Name":"joe"},{"Name":"sally"}]
//  }
package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Resource interface{}

type getter interface {
	Get(values url.Values) (int, interface{})
}
type poster interface {
	Post(values url.Values) (int, interface{})
}
type putter interface {
	Put(values url.Values) (int, interface{})
}
type deleter interface {
	Delete(values url.Values) (int, interface{})
}

type BaseResource struct{}

func (BaseResource) Get(values url.Values) (int, interface{}) {
	return 405, ""
}

func (BaseResource) Post(values url.Values) (int, interface{}) {
	return 405, ""
}

func (BaseResource) Put(values url.Values) (int, interface{}) {
	return 405, ""
}

func (BaseResource) Delete(values url.Values) (int, interface{}) {
	return 405, ""
}

type API struct{}

func (api *API) Abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
}

func (api *API) requestHandler(resource Resource) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		var (
			data interface{}
			code int
		)

		request.ParseForm()
		method, values := request.Method, request.Form

		if resource, ok := resource.(getter); ok && method == GET {
			code, data = resource.Get(values)
		} else if resource, ok := resource.(poster); ok && method == POST {
			code, data = resource.Post(values)
		} else if resource, ok := resource.(putter); ok && method == PUT {
			code, data = resource.Put(values)
		} else if resource, ok := resource.(deleter); ok && method == DELETE {
			code, data = resource.Delete(values)
		} else {
			api.Abort(rw, 405)
			return
		}

		responseWriter := json.NewEncoder(rw)
		rw.WriteHeader(code)
		if responseWriter.Encode(data) != nil {
			api.Abort(rw, 500)
			return
		}
	}
}

func (api *API) AddResource(path string, resource Resource) {
	http.HandleFunc(path, api.requestHandler(resource))
}

func (api *API) Start(port int) {
	portString := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portString, nil)
}
