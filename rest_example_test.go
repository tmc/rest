package rest_test

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

func ExampleAPI_Start() {
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
