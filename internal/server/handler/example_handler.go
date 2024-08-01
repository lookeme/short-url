package handler

// ExamplePost saves a URL sent in the request body to the database and returns a success response.
// The function accepts an HTTP response writer and request as arguments.
// The request should have a POST method on the root URL ("/") with a payload in JSON format.
// The payload should contain a "url" field with the URL to be saved.
// The function returns an HTTP 200 OK response with a "Content-Type" of "application/json" and a status of "OK".
// Example usage:
//
//	req, _ := http.NewRequest("POST", "/", bytes.NewReader([]byte(`{"url":"www.example.ru"}`)))
//	res := httptest.NewRecorder()
//	ExamplePost(res, req)
//	fmt.Println(res.Body.String())
//
// Output:
//
//	{"message":"URL saved successfully"}
func ExamplePost() {
	// URLHandler{}.HandlePOST(res http.ResponseWriter, req *http.Request) // Save URL
	// Output:
	// # Request
	// POST /
	//  Content-Type: application/json
	// Body { "url" : "www.example.ru" }
	// # Response
	// HTTP/1.1 200 OK
	// Content-Type: application/json
	// Status: OK

}
