package response

import "net/http"

func HxRedirect(rw http.ResponseWriter, r *http.Request) {
	res := new(rw, r)

	res.statusCode = http.StatusOK
	headers := make(map[string]string)
	headers["HX-Redirect"] = "/"
	headers["Content-Type"] = "text/html; charset=utf-8"
	res.headers = headers

	res.write()
}
