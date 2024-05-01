// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.648
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"errors"
	"net/http"
	"strconv"

	ol "github.com/kencx/dusk/integrations/openlibrary"
	"github.com/kencx/dusk/ui/partials"
	"github.com/kencx/dusk/ui/shared"
)

var ErrNotValidIsbn = errors.New("invalid isbn")

type Search struct {
	DefaultTab string
	Results    ol.QueryResults
	Message    string
	Err        error
}

func ImportRenderResults(rw http.ResponseWriter, r *http.Request, res ol.QueryResults) {
	ImportResults(res, "", nil).Render(r.Context(), rw)
}

func ImportResultsMessage(rw http.ResponseWriter, r *http.Request, res ol.QueryResults, message string) {
	ImportResults(res, message, nil).Render(r.Context(), rw)
}

func ImportResultsError(rw http.ResponseWriter, r *http.Request, err error) {
	ImportResults(nil, "", err).Render(r.Context(), rw)
}

func (v Search) Render(rw http.ResponseWriter, r *http.Request) {
	v.Html().Render(r.Context(), rw)
}

func (v *Search) Html() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h2>Add Books</h2>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = partials.Tabs(ImportTabs, v.DefaultTab).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" <div class=\"import__result_list\"></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = shared.Base().Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func searchForm() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form hx-post=\"/search\" hx-target=\".import__result_list\" hx-swap=\"innerHTML\" hx-indicator=\".spinner\"><fieldset role=\"group\"><input id=\"search\" name=\"openlibrary\" placeholder=\"Search for an ISBN, title or author\"> <button type=\"submit\">Submit</button></fieldset><small><a href=\"https://www.isbn-13.info/example\">ISBNs</a> must contain 10 or 13 characters, excluding dashes and spaces.</small><div class=\"spinner\" aria-busy=\"true\"></div></form>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func ImportResults(results ol.QueryResults, message string, err error) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var4 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var4 == nil {
			templ_7745c5c3_Var4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if err == ErrNotValidIsbn {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"card error fluid\"><p>You entered an invalid ISBN. A book's ISBN is usually found on the back cover, near the barcode. It will contain 10 or 13 characters, plus any dashes and spaces.</p><p>Valid examples: 978-0495011606, 9780136006176, 0077354761, 013603599X</p><p>Alternatively, you can try searching for the book's title or author.</p></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else if err != nil {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"card error fluid\"><p>Something went wrong, please try again</p></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if message != "" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"card fluid message\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templ.Raw(message).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if err == nil {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h2>Results <small>Found ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(len(results)))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/search.templ`, Line: 95, Col: 44}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" result(s)</small></h2><form hx-indicator=\".spinner\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, result := range results {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"import__result\"><img alt=\"\" src=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var6 string
				templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(result.CoverUrl)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/search.templ`, Line: 100, Col: 38}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><div class=\"details\"><h4>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var7 string
				templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(result.Title)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/search.templ`, Line: 103, Col: 21}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" <small>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				for _, author := range result.Authors {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span class=\"author\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var8 string
					templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(author)
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/search.templ`, Line: 106, Col: 38}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</small></h4><ul><li>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if len(result.Isbn10) > 0 || len(result.Isbn13) > 0 {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("ISBN: ")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					if len(result.Isbn10) > 0 {
						var templ_7745c5c3_Var9 string
						templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(result.Isbn10[0])
						if templ_7745c5c3_Err != nil {
							return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/search.templ`, Line: 115, Col: 28}
						}
						_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" ")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					if len(result.Isbn13) > 0 {
						var templ_7745c5c3_Var10 string
						templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(result.Isbn13[0])
						if templ_7745c5c3_Err != nil {
							return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/search.templ`, Line: 118, Col: 28}
						}
						_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
						if templ_7745c5c3_Err != nil {
							return templ_7745c5c3_Err
						}
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li><li>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if result.PublishDate != "" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("Published: ")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var11 string
					templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(result.PublishDate)
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/search.templ`, Line: 124, Col: 40}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li></ul></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if len(result.Isbn10) > 0 {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<input type=\"hidden\" name=\"result\" value=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var12 string
					templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(result.Isbn10[0])
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/search.templ`, Line: 130, Col: 65}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> ")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				if len(result.Isbn13) > 0 {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<input type=\"hidden\" name=\"result\" value=\"")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var13 string
					templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(result.Isbn13[0])
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/search.templ`, Line: 133, Col: 65}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"actions\"><select name=\"tag-option\" hx-trigger=\"change\" hx-post=\"/search/add\" hx-target=\".import__result_list\" hx-swap=\"innerHTML\" hx-include=\"this\"><option value=\"add\">Add book</option> <option value=\"to-read\">To read</option> <option value=\"reading\">Reading</option> <option value=\"read\">Read</option></select></div></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"spinner\"></div></form>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}