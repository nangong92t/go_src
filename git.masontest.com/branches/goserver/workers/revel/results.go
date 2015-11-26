package revel

import (
//	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"reflect"
	"time"
)

type Result interface {
	Apply(req *Request, resp *Response) interface{}
    Clean()
}

// This result handles all kinds of error codes (500, 404, ..).
// It renders the relevant error page (errors/CODE.format, e.g. errors/500.json).
// If RunMode is "dev", this results in a friendly error page.
type ErrorResult struct {
	RenderArgs map[string]interface{}
	Error      error
}

func (r ErrorResult) Apply(req *Request, resp *Response) interface{} {
	// If it's not a revel error, wrap it in one.
    panic(r.Error)

    return nil
}

func (r ErrorResult) Clean() {}

type PlaintextErrorResult struct {
	Error error
}

// This method is used when the template loader or error template is not available.
func (r PlaintextErrorResult) Apply(req *Request, resp *Response) interface{} {
    panic(r.Error)
    return nil
}
func (r PlaintextErrorResult) Clean() {}

// Action methods return this result to request a template be rendered.
type RenderTemplateResult struct {
	Template   Template
	RenderArgs map[string]interface{}
}

func (r *RenderTemplateResult) Apply(req *Request, resp *Response) interface{} {
	// If it's a HEAD request, throw away the bytes.
	return resp.Out
}
func (r *RenderTemplateResult) Clean() {}

func (r *RenderTemplateResult) render(req *Request, resp *Response, wr io.Writer) interface{} {
    return r.RenderArgs
}

type RenderHtmlResult struct {
	html string
}

func (r RenderHtmlResult) Clean() {}
func (r RenderHtmlResult) Apply(req *Request, resp *Response) interface{} {
	return r.html
}

type RenderRpcResult struct {
    Error       error
    ErrorCode   int
    Data        interface{}
}

func (r RenderRpcResult) Apply(req *Request, resp *Response) interface{} {
    return r
}

func (r RenderRpcResult) Clean() {
    r.Error     = nil
    r.ErrorCode = 0
    r.Data      = nil
}

type RenderJsonResult struct {
	obj      interface{}
	callback string
}

func (r RenderJsonResult) Clean() {}

func (r RenderJsonResult) Apply(req *Request, resp *Response) interface{} {
    return r.obj
}

type RenderXmlResult struct {
	obj interface{}
}

func (r RenderXmlResult) Clean() {}
func (r RenderXmlResult) Apply(req *Request, resp *Response) interface{} {
	var b []byte
	var err error
	if Config.BoolDefault("results.pretty", false) {
		b, err = xml.MarshalIndent(r.obj, "", "  ")
	} else {
		b, err = xml.Marshal(r.obj)
	}

	if err != nil {
	    return ErrorResult{Error: err}.Apply(req, resp)
	}

    return b
}

type RenderTextResult struct {
	text string
}

func (r RenderTextResult) Clean() {}
func (r RenderTextResult) Apply(req *Request, resp *Response) interface{} {
	return r.text
}

type ContentDisposition string

var (
	Attachment ContentDisposition = "attachment"
	Inline     ContentDisposition = "inline"
)

type BinaryResult struct {
	Reader   io.Reader
	Name     string
	Length   int64
	Delivery ContentDisposition
	ModTime  time.Time
}

func (r *BinaryResult) Clean() {}

func (r *BinaryResult) Apply(req *Request, resp *Response) interface{} {
    return r.Reader
}

type RedirectToUrlResult struct {
	url string
}

func (r *RedirectToUrlResult) Clean() {}

func (r *RedirectToUrlResult) Apply(req *Request, resp *Response) interface{} {
	resp.Out.Header().Set("Location", r.url)
	resp.WriteHeader(http.StatusFound, "")

    return nil
}

type RedirectToActionResult struct {
	val interface{}
}

func (r *RedirectToActionResult) Clean() {}

func (r *RedirectToActionResult) Apply(req *Request, resp *Response) interface{} {
	url, err := getRedirectUrl(r.val)
	if err != nil {
		ERROR.Println("Couldn't resolve redirect:", err.Error())
		return ErrorResult{Error: err}.Apply(req, resp)
	}
    return url
}

func getRedirectUrl(item interface{}) (string, error) {
	// Handle strings
	if url, ok := item.(string); ok {
		return url, nil
	}

	// Handle funcs
	val := reflect.ValueOf(item)
	typ := reflect.TypeOf(item)
	if typ.Kind() == reflect.Func && typ.NumIn() > 0 {
		// Get the Controller Method
		recvType := typ.In(0)
		method := FindMethod(recvType, val)
		if method == nil {
			return "", errors.New("couldn't find method")
		}

		// Construct the action string (e.g. "Controller.Method")
		if recvType.Kind() == reflect.Ptr {
			recvType = recvType.Elem()
		}
		action := recvType.Name() + "." + method.Name
		actionDef := MainRouter.Reverse(action, make(map[string]string))
		if actionDef == nil {
			return "", errors.New("no route for action " + action)
		}

		return actionDef.String(), nil
	}

	// Out of guesses
	return "", errors.New("didn't recognize type: " + typ.String())
}
