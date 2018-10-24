package lazy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/gorilla/mux"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()

type Response struct {
	Error string      `json:"error,omitempt"`
	Data  interface{} `json:"data,omitempt"`
}

// Router routes all request to rest API services.
type Router struct {
	router *mux.Router
}

type endpoint struct {
	service     interface{}
	serviceType reflect.Type
	dataType    reflect.Type

	get    reflect.Method
	put    reflect.Method
	new    reflect.Method
	delete reflect.Method
}

// isExported() and isExportedOrBuiltinType() from net/rpc
// Copyright 2009 The Go Authors governed by a BSD-style
// license

// Is this an exported - upper case - name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

func sendJsonResponse(w http.ResponseWriter, data interface{}) {
	resp := Response{
		Data: data,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Marshal error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (e *endpoint) findGet() error {
	getMethod, ok := e.serviceType.MethodByName("Get")
	if !ok {
		return fmt.Errorf("Service does not have a Get Method")
	}

	t := getMethod.Type

	if t.NumIn() != 3 {
		return fmt.Errorf("Get method need 3 arguments, has %d", t.NumIn())
	}

	if t.In(1) != typeOfContext {
		return fmt.Errorf("Ctx argument must be of type context.Context.  Found %v instead", t.In(1))
	}

	if t.In(2).Kind() != reflect.Int {
		return fmt.Errorf("Id argument must be an int.  Found %v instead", t.In(2))
	}

	if t.NumOut() != 2 {
		return fmt.Errorf("Get method needs 2 return values, has %d", t.NumOut())
	}

	e.dataType = t.Out(0)
	if e.dataType.Kind() != reflect.Ptr {
		return fmt.Errorf("Data type %v must be a pointer", t.Out(0))
	}

	if !isExportedOrBuiltinType(e.dataType) {
		return fmt.Errorf("Data type %v not exported", t.Out(0))
	}

	if t.Out(1) != typeOfError {
		return fmt.Errorf("Second return type must be error; Found %v instead", t.Out(1))
	}

	e.get = getMethod
	return nil
}

func (e *endpoint) findPut() error {
	putMethod, ok := e.serviceType.MethodByName("Put")
	if !ok {
		return fmt.Errorf("Service does not have a Put Method")
	}

	t := putMethod.Type

	if t.NumIn() != 4 {
		return fmt.Errorf("Put method needs 4 arguments, has %d", t.NumIn())
	}

	if t.In(1) != typeOfContext {
		return fmt.Errorf("Ctx argument must be of type context.Context.  Found %v instead", t.In(1))
	}

	if t.In(2).Kind() != reflect.Int {
		return fmt.Errorf("Id argument mus be an int.  Found %v instead", t.In(2))
	}

	if t.In(3) != e.dataType {
		return fmt.Errorf("data argument mus be %v.  Found %v instead", e.dataType, t.In(3))
	}

	if t.NumOut() != 1 {
		return fmt.Errorf("Get method needs 1 return value, has %d", t.NumOut())
	}

	if t.Out(0) != typeOfError {
		return fmt.Errorf("First return type must be error; Found %v instead", t.Out(0))
	}

	e.put = putMethod
	return nil
}

func (e *endpoint) findNew() error {
	newMethod, ok := e.serviceType.MethodByName("New")
	if !ok {
		return fmt.Errorf("Service does not have a New Method")
	}

	t := newMethod.Type

	if t.NumIn() != 3 {
		return fmt.Errorf("New method needs 3 arguments, has %d", t.NumIn())
	}

	if t.In(1) != typeOfContext {
		return fmt.Errorf("Ctx argument must be of type context.Context.  Found %v instead", t.In(1))
	}

	if t.In(2) != e.dataType {
		return fmt.Errorf("data argument mus be %v.  Found %v instead", e.dataType, t.In(2))
	}

	if t.NumOut() != 2 {
		return fmt.Errorf("Get method needs 1 return value, has %d", t.NumOut())
	}

	if t.Out(0).Kind() != reflect.Int {
		return fmt.Errorf("First return type must be an int.  Found %v instead", t.Out(0))
	}

	if t.Out(1) != typeOfError {
		return fmt.Errorf("Second return type must be error; Found %v instead", t.Out(1))
	}

	e.new = newMethod
	return nil
}

func (e *endpoint) findDelete() error {
	deleteMethod, ok := e.serviceType.MethodByName("Delete")
	if !ok {
		return fmt.Errorf("Service does not have a Delete Method")
	}

	t := deleteMethod.Type

	if t.NumIn() != 3 {
		return fmt.Errorf("Delete method needs 3 arguments, has %d", t.NumIn())
	}

	if t.In(1) != typeOfContext {
		return fmt.Errorf("Ctx argument must be of type context.Context.  Found %v instead", t.In(1))
	}

	if t.In(2).Kind() != reflect.Int {
		return fmt.Errorf("Id argument mus be an int.  Found %v instead", t.In(2))
	}

	if t.NumOut() != 1 {
		return fmt.Errorf("Delete method needs 1 return value, has %d", t.NumOut())
	}

	if t.Out(0) != typeOfError {
		return fmt.Errorf("First return type must be error; Found %v instead", t.Out(0))
	}

	e.delete = deleteMethod
	return nil
}

func handleCallError(method string, v reflect.Value, w http.ResponseWriter) bool {
	if !v.IsNil() {
		err := v.Interface().(error)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}

	return false
}

func (e *endpoint) decodeData(r *http.Request) (*reflect.Value, error) {
	data := reflect.New(e.dataType.Elem())
	err := json.NewDecoder(r.Body).Decode(data.Interface())
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (e *endpoint) handleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	values := e.get.Func.Call([]reflect.Value{
		reflect.ValueOf(e.service),
		reflect.ValueOf(r.Context()),
		reflect.ValueOf(id)})
	data := values[0].Interface()

	if handleCallError("Get", values[1], w) {
		return
	}

	sendJsonResponse(w, data)
}

func (e *endpoint) handlePut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	data, err := e.decodeData(r)
	if err != nil {
		log.Printf("Decode error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	values := e.put.Func.Call([]reflect.Value{
		reflect.ValueOf(e.service),
		reflect.ValueOf(r.Context()),
		reflect.ValueOf(id),
		*data})
	if handleCallError("Put", values[0], w) {
		return
	}

	sendJsonResponse(w, id)
}

func (e *endpoint) handleNew(w http.ResponseWriter, r *http.Request) {
	data, err := e.decodeData(r)
	if err != nil {
		log.Printf("Decode error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	values := e.new.Func.Call([]reflect.Value{
		reflect.ValueOf(e.service),
		reflect.ValueOf(r.Context()),
		*data})
	if handleCallError("New", values[1], w) {
		return
	}

	sendJsonResponse(w, values[0].Interface())
}

func (e *endpoint) handleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	values := e.delete.Func.Call([]reflect.Value{
		reflect.ValueOf(e.service),
		reflect.ValueOf(r.Context()),
		reflect.ValueOf(id)})
	if handleCallError("Delete", values[0], w) {
		return
	}

	sendJsonResponse(w, id)
}

// NewRouter creates a new Router.
func NewRouter() *Router {
	return &Router{
		router: mux.NewRouter(),
	}
}

// AddService adds a service to the router.
func (r *Router) AddService(prefix string, service interface{}) error {
	e := &endpoint{
		service:     service,
		serviceType: reflect.TypeOf(service),
	}

	// TODO(konkers): Support partial endpoints
	err := e.findGet()
	if err != nil {
		return err
	}
	err = e.findPut()
	if err != nil {
		return err
	}
	err = e.findNew()
	if err != nil {
		return err
	}
	err = e.findDelete()
	if err != nil {
		return err
	}

	s := r.router.PathPrefix("/" + prefix).Subrouter()
	s.HandleFunc("/get/{id:[0-9]+}", e.handleGet)
	s.HandleFunc("/put/{id:[0-9]+}", e.handlePut)
	s.HandleFunc("/new", e.handleNew)
	s.HandleFunc("/delete/{id:[0-9]+}", e.handleDelete)

	return nil
}

// ServeHTTP implements the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
