package handler

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
