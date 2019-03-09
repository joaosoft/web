package main

import (
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

	requestGetBoundary(c)

	requestAuthBasic(c)
	requestAuthJwt(c)

	requestOptionsOK(c)
	requestOptionsNotFound(c)

	bindFormData(c)
}

func requestGet(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/hello/joao?a=1&b=2&c=1,2,3")
	if err != nil {
		panic(err)
	}

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", response)
}

func requestGetBoundary(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/hello/joao/download")
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
	request, err := c.NewRequest(web.MethodOptions, "localhost:9001/auth-basic")
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
	request, err := c.NewRequest(web.MethodOptions, "localhost:9001/auth-basic-invalid")
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
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/auth-basic")
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
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/auth-jwt")
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
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/form-data")
	if err != nil {
		panic(err)
	}

	request.SetFormData("var_one", "one")
	request.SetFormData("var_two", "2")

	response, err := request.WithContentType(web.ContentTypeMultipartFormData).Send()
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

	fmt.Printf("\n\n%+v", formData)
}
