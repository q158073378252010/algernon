// Package httpclient provides Lua functions for a HTTP client
package httpclient

import (
	"github.com/ddliu/go-httpclient"
	log "github.com/sirupsen/logrus"
	"github.com/xyproto/gopher-lua"
)

const (
	// HTTPClientClass is an identifier for the HTTPClient class in Lua
	HTTPClientClass = "HTTPClient"
)

// Get the first argument, "self", and cast it from userdata to a library (which is really a hash map).
func checkHTTPClientClass(L *lua.LState) *httpclient.HttpClient {
	ud := L.CheckUserData(1)
	if hc, ok := ud.Value.(*httpclient.HttpClient); ok {
		return hc
	}
	L.ArgError(1, "HTTPClient expected")
	return nil
}

// Create a new httpclient.HttpClient.
// The first argument is the language code, used with Accept-Language,
// and is optional.
func constructHTTPClient(L *lua.LState, userAgent string) (*lua.LUserData, error) {
	// Use the first argument as the name of the tag
	language := L.ToString(1)
	if language == "" {
		language = "en-us"
	}

	// Create a new HTTP Client
	hc := httpclient.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: userAgent,
		"Accept-Language":        language,
	})

	// Create a new userdata struct
	ud := L.NewUserData()
	ud.Value = hc
	L.SetMetatable(ud, L.GetTypeMetatable(HTTPClientClass))
	return ud, nil
}

// Use the http client to GET the given URL.
func hcGet(L *lua.LState) int {
	hc := checkHTTPClientClass(L) // arg 1
	URL := L.ToString(2)
	if URL == "" {
		L.ArgError(2, "URL expected")
		return 0 // no results
	}

	urlValues := L.ToTable(3)
	if urlValues == nil {
		L.ArgError(3, "keys and values from the URL is expected, use {} to skip")
		return 0 // no results
	}

	keysAndValues, _, _, _ := Table2maps(urlValues)

	println("GET " + URL)

	keysAndValues := map[string]string{
		"q": "news",
	}

	// Fetch the given URL
	res, err := hc.Get(URL, keysAndValues)

	println(res.StatusCode, err)

	// Return a string
	contents := "OLLA BOLLA"
	L.Push(lua.LString(contents))
	return 1 // number of results
}

func hcString(L *lua.LState) int {
	L.Push(lua.LString("HTTP client based on github.com/ddliu/go-httpclient"))
	return 1 // number of results
}

// The hash map methods that are to be registered
var hcMethods = map[string]lua.LGFunction{
	"__tostring":   hcString,
	"Get":          hcGet,
	"Post":         hcPost,
	"PutJSON":      hcPutJSON,
	"Delete":       hcDelete,
	"Options":      hcOptions,
	"Head":         hcHead,
	"SetUserAgent": hcSetUserAgent,
	"SetLanguage":  hcSetLanguage,
	"SetHeader":    hcSetHeader,
	"SetOption":    hcSetOption,
	"SetCookie":    hcSetCookie,
}

// Load makes functions related to httpclient available to the given Lua state
func Load(L *lua.LState, userAgent string) {

	println("LOAD HTTP CLIENT")

	// Register the HTTPClient class and the methods that belongs with it.
	metaTableHC := L.NewTypeMetatable(HTTPClientClass)
	metaTableHC.RawSetH(lua.LString("__index"), metaTableHC)
	L.SetFuncs(metaTableHC, hcMethods)

	// The constructor for HTTPClient
	L.SetGlobal("HTTPClient", L.NewFunction(func(L *lua.LState) int {
		// Construct a new HTTPClient
		userdata, err := constructHTTPClient(L, userAgent)
		if err != nil {
			log.Error(err)
			L.Push(lua.LString(err.Error()))
			return 1 // Number of returned values
		}

		// Return the Lua Page object
		L.Push(userdata)
		return 1 // number of results
	}))

	L.SetGlobal("httpGet", L.NewFunction(func(L *lua.LState) int {
		httpclient.Defaults(httpclient.Map{
			httpclient.OPT_USERAGENT: userAgent,
			"Accept-Language":        "en-us",
		})
		retval := "OSTEBOLLE"

		L.Push(retval)
		return 1 // number of results
	}))
}
