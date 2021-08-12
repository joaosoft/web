package main

import (
	"encoding/json"
	"fmt"
	"web"

	"github.com/joaosoft/auth-types/jwt"
)

func main() {
	// create a new client
	c, err := web.NewClient(web.WithClientMultiAttachmentMode(web.MultiAttachmentModeBoundary))
	if err != nil {
		panic(err)
	}

	requestGet(c)
	requestPost(c)

	requestGetBoundary(c)

	requestAuthBasic(c)
	requestAuthJwt(c)

	requestOptionsOK(c)
	requestOptionsNotFound(c)

	bindFormData(c)
	bindUrlFormData(c)
}

func requestGet(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/hello/joao?a=1&b=2&c=1,2,3", web.ContentTypeApplicationJSON, nil)
	if err != nil {
		panic(err)
	}

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", response)
}

func requestPost(c *web.Client) {
	request, err := c.NewRequest(web.MethodPost, "localhost:9001/hello/joao?a=1&b=2&c=1,2,3", web.ContentTypeApplicationJSON, nil)
	if err != nil {
		panic(err)
	}

	data := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "joao",
		Age:  30,
	}

	bytes, _ := json.Marshal(data)

	response, err := request.WithBody(bytes).Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", string(response.Body))
}

func requestGetBoundary(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/hello/joao/download", web.ContentTypeApplicationJSON, nil)
	if err != nil {
		panic(err)
	}

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", response)
}

func requestOptionsOK(c *web.Client) {
	request, err := c.NewRequest(web.MethodOptions, "localhost:9001/auth-basic", web.ContentTypeApplicationJSON, nil)
	if err != nil {
		panic(err)
	}

	_, err = request.WithAuthBasic("joao", "ribeiro")
	if err != nil {
		panic(err)
	}

	request.SetHeader(web.HeaderAccessControlRequestMethod, []string{string(web.MethodGet)})
	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n%d: %s\n\n", response.Status, string(response.Body))
}

func requestOptionsNotFound(c *web.Client) {
	request, err := c.NewRequest(web.MethodOptions, "localhost:9001/auth-basic-invalid", web.ContentTypeApplicationJSON, nil)
	if err != nil {
		panic(err)
	}

	_, err = request.WithAuthBasic("joao", "ribeiro")
	if err != nil {
		panic(err)
	}

	request.SetHeader(web.HeaderAccessControlRequestMethod, []string{string(web.MethodGet)})
	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n%d: %s\n\n", response.Status, string(response.Body))
}

func requestAuthBasic(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/auth-basic", web.ContentTypeApplicationJSON, nil)
	if err != nil {
		panic(err)
	}

	_, err = request.WithAuthBasic("joao", "ribeiro")
	if err != nil {
		panic(err)
	}

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n%d: %s\n\n", response.Status, string(response.Body))
}

func requestAuthJwt(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/auth-jwt", web.ContentTypeApplicationJSON, nil)
	if err != nil {
		panic(err)
	}

	_, err = request.WithAuthJwt(jwt.Claims{"name": "joao", "age": 30}, "bananas")
	if err != nil {
		panic(err)
	}

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n%d: %s\n\n", response.Status, string(response.Body))
}

func bindFormData(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/form-data", web.ContentTypeMultipartFormData, nil)
	if err != nil {
		panic(err)
	}

	request.SetFormData("var_one", "one")
	request.SetFormData("var_two", "2")

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	formData := struct {
		VarOne string `json:"var_one"`
		VarTwo int    `json:"var_two"`
	}{}

	if err := response.BindFormData(&formData); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\nvar_one: %s", response.GetFormDataString("var_one"))
	fmt.Printf("\nvar_two: %s", response.GetFormDataString("var_two"))

	fmt.Printf("\n\nFORM DATA: %+v\n", formData)
}

func bindUrlFormData(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/url-form-data", web.ContentTypeApplicationForm, nil)
	if err != nil {
		panic(err)
	}

	request.SetFormData("var_one", "one")
	request.SetFormData("var_two", "2")

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	formData := struct {
		VarOne string `json:"var_one"`
		VarTwo int    `json:"var_two"`
	}{}

	if err := response.BindFormData(&formData); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\nvar_one: %s", response.GetFormDataString("var_one"))
	fmt.Printf("\nvar_two: %s", response.GetFormDataString("var_two"))

	fmt.Printf("\n\nURL FORM DATA: %+v\n", formData)
}
