// Package httpclient provides Lua functions for a HTTP client
package httpclient

import (
	"github.com/ddliu/go-httpclient"
	log "github.com/sirupsen/logrus"
	"github.com/xyproto/algernon/lua/convert"
	"github.com/xyproto/gopher-lua"
	"io/ioutil"
	"strings"
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

// Create a new httpclient.HttpClient. The Lua function takes no arguments.
func constructHTTPClient(L *lua.LState, userAgent string) (*lua.LUserData, error) {
	// Create a new HTTP Client
	hc := httpclient.NewHttpClient()

	// Set the default user agent to the server name
	hc.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: userAgent,
	})

	// Create a new userdata struct
	ud := L.NewUserData()
	ud.Value = hc
	L.SetMetatable(ud, L.GetTypeMetatable(HTTPClientClass))
	return ud, nil
}

// hcGet is a Lua function for running the GET method on a given URL.
// It also takes a table with URL arguments (optional).
func hcGet(L *lua.LState) int {
	hc := checkHTTPClientClass(L) // arg 1
	URL := L.ToString(2)          // arg 2
	if URL == "" {
		L.ArgError(2, "URL expected")
		return 0 // no results
	}

	// Request headers
	headers := make(map[string]string)
	urlValues := L.ToTable(3)
	if urlValues != nil {
		headers, _, _, _ = convert.Table2maps(urlValues)
	}

	// GET the given URL
	resp, err := hc.Do("GET", URL, headers, nil)
	if err != nil {
		log.Error(err)
		return 0 // no results
	}

	// Read the returned body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return 0 // no results
	}

	// Return a string
	L.Push(lua.LString(string(body)))
	return 1 // number of results
}

// hcHead is a Lua function for running the HEAD method on a given URL.
// It also takes a table with URL arguments (optional).
func hcHead(L *lua.LState) int {
	hc := checkHTTPClientClass(L) // arg 1
	URL := L.ToString(2)          // arg 2
	if URL == "" {
		L.ArgError(2, "URL expected")
		return 0 // no results
	}

	// Request headers
	headers := make(map[string]string)
	urlValues := L.ToTable(3)
	if urlValues != nil {
		headers, _, _, _ = convert.Table2maps(urlValues)
	}

	// HEAD the given URL
	resp, err := hc.Do("HEAD", URL, headers, nil)
	if err != nil {
		log.Error(err)
		return 0 // no results
	}

	// Read the returned body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return 0 // no results
	}

	// Return a string
	L.Push(lua.LString(string(body)))
	return 1 // number of results
}

// hcDelete is a Lua function for running the DELETE method on a given URL.
// It also takes a table with URL arguments (optional).
func hcDelete(L *lua.LState) int {
	hc := checkHTTPClientClass(L) // arg 1
	URL := L.ToString(2)          // arg 2
	if URL == "" {
		L.ArgError(2, "URL expected")
		return 0 // no results
	}

	// Request headers
	headers := make(map[string]string)
	urlValues := L.ToTable(3)
	if urlValues != nil {
		headers, _, _, _ = convert.Table2maps(urlValues)
	}

	// DELETE the given URL
	resp, err := hc.Do("DELETE", URL, headers, nil)
	if err != nil {
		log.Error(err)
		return 0 // no results
	}

	// Read the returned body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return 0 // no results
	}

	// Return a string
	L.Push(lua.LString(string(body)))
	return 1 // number of results
}

// hcPost is a Lua function for running the POST method on a given URL.
// It also takes a table with URL arguments (optional).
// It can also takes a string to post as the body (optional).
func hcPost(L *lua.LState) int {
	hc := checkHTTPClientClass(L) // arg 1
	URL := L.ToString(2)          // arg 2
	if URL == "" {
		L.ArgError(2, "URL expected")
		return 0 // no results
	}

	// Request headers
	headers := make(map[string]string)
	urlValues := L.ToTable(3) // arg 3 (optional)
	if urlValues != nil {
		headers, _, _, _ = convert.Table2maps(urlValues)
	}

	// Body
	body := L.ToString(4) // arg 4 (optional)

	// POST the given URL
	resp, err := hc.Do("POST", URL, headers, strings.NewReader(body))
	if err != nil {
		log.Error(err)
		return 0 // no results
	}

	// Read the returned body
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return 0 // no results
	}

	// Return a string
	L.Push(lua.LString(string(respBody)))
	return 1 // number of results
}

// hcString is a Lua function that returns a descriptive string
func hcString(L *lua.LState) int {
	L.Push(lua.LString("HTTP client based on github.com/ddliu/go-httpclient"))
	return 1 // number of results
}

// hcSetUserAgent is a Lua function for setting the user agent string
func hcSetUserAgent(L *lua.LState) int {
	hc := checkHTTPClientClass(L) // arg 1
	userAgent := L.ToString(2)    // arg 2
	if userAgent == "" {
		L.ArgError(2, "User agent string expected")
		return 0 // no results
	}

	hc.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: userAgent,
	})

	return 0 // no results
}

// The hash map methods that are to be registered
var hcMethods = map[string]lua.LGFunction{
	"__tostring":   hcString,
	"SetUserAgent": hcSetUserAgent,
	"Get":          hcGet,
	"Head":         hcHead,
	"Delete":       hcDelete,
	"Post":         hcPost,
	//"PutJSON":      hcPutJSON,
	//"Options":      hcOptions,
	//"SetLanguage":  hcSetLanguage,
	//"SetHeader":    hcSetHeader,
	//"SetOption":    hcSetOption,
	//"SetCookie":    hcSetCookie,
}

// Load makes functions related to httpclient available to the given Lua state
func Load(L *lua.LState, userAgent string) {

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
}
