package response

import "net/http"

func HxRedirect(rw http.ResponseWriter, r *http.Request, path string) {
	res := new(rw, r)

	res.statusCode = http.StatusOK
	headers := make(map[string]string)
	headers["HX-Redirect"] = path
	headers["Content-Type"] = "text/html; charset=utf-8"
	res.headers = headers

	res.write()
}
