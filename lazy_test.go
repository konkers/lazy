package lazy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/phayes/freeport"
)

type TestData struct {
	ID   int
	Name string
}

type TestService struct {
	nextID int
	data   map[int]*TestData

	FailNew bool
}

func (s *TestService) Get(id int) (*TestData, error) {
	data, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("ID %d does not exist.", id)
	}
	return data, nil
}

func (s *TestService) Put(id int, data *TestData) error {
	_, ok := s.data[id]
	if !ok {
		return fmt.Errorf("ID %d does not exist", id)
	}
	data.ID = id
	s.data[id] = data
	return nil
}

func (s *TestService) New(data *TestData) (int, error) {
	if s.FailNew {
		return 0, fmt.Errorf("New Failure")
	}

	data.ID = s.nextID
	s.nextID = s.nextID + 1
	s.data[data.ID] = data

	return data.ID, nil
}

func (s *TestService) Delete(id int) error {
	_, ok := s.data[id]
	if !ok {
		return fmt.Errorf("ID %d does not exist", id)
	}

	delete(s.data, id)
	return nil
}

func NewTestService() *TestService {
	return &TestService{
		data:   make(map[int]*TestData),
		nextID: 1,
	}
}

type DummyResponseWriter struct {
	header     http.Header
	StatusCode int
}

func NewDummyResponseWriter() *DummyResponseWriter {
	return &DummyResponseWriter{
		header:     make(http.Header),
		StatusCode: http.StatusOK,
	}
}

func (w *DummyResponseWriter) Header() http.Header {
	return w.header
}

func (w *DummyResponseWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func (w *DummyResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func testNewReq(addr string, data *TestData) (int, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		return 0, fmt.Errorf("Can't encode json: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/test/new", addr),
		"application/json; charset=utf-8", buf)
	if err != nil {
		return 0, fmt.Errorf("Post error: %v", err)
	}

	var ret struct {
		Error string `json:"error,omitempt"`
		Data  int    `json:"data,omitempt"`
	}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return 0, fmt.Errorf("JSON decode error: %v", err)
	}

	return ret.Data, nil
}

func testAndValidateNewReq(t *testing.T, addr string, data *TestData) {
	id, err := testNewReq(addr, data)
	if err != nil {
		t.Errorf("Can't create data: %v", err)
	}
	if id != data.ID {
		t.Errorf("Expected id of %d, got %d instead.", data.ID, id)
	}
}

func testGetReq(addr string, id int) (*TestData, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/test/get/%d", addr, id))
	if err != nil {
		return nil, fmt.Errorf("Can't get: %v", err)
	}

	var ret struct {
		Error string    `json:"error,omitempt"`
		Data  *TestData `json:"data,omitempt"`
	}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return nil, fmt.Errorf("JSON decode error: %v", err)
	}

	return ret.Data, nil
}

func testAndValidateGetReq(t *testing.T, addr string, id int, name string) {
	data, err := testGetReq(addr, id)
	if err != nil {
		t.Errorf("Can't get %d: %v", id, err)
	}
	if data.ID != id {
		t.Errorf("Expected id %d, got %d instead.", id, data.ID)
	}
	if data.Name != name {
		t.Errorf("Expected name %s, got %s instead.", name, data.Name)
	}
}

func testPutReq(addr string, id int, data *TestData) (int, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		return 0, fmt.Errorf("Can't encode json: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/test/put/%d", addr, id),
		"application/json; charset=utf-8", buf)
	if err != nil {
		return 0, fmt.Errorf("Post error: %v", err)
	}

	var ret struct {
		Error string `json:"error,omitempt"`
		Data  int    `json:"data,omitempt"`
	}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return 0, fmt.Errorf("JSON decode error: %v", err)
	}

	return ret.Data, nil
}

func testAndValidatePutReq(t *testing.T, addr string, data *TestData) {
	id, err := testPutReq(addr, data.ID, data)
	if err != nil {
		t.Errorf("Can't create data: %v", err)
	}
	if id != data.ID {
		t.Errorf("Expected id of %d, got %d instead.", data.ID, id)
	}
}

func testDeleteReq(addr string, id int) (int, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/test/delete/%d", addr, id))
	if err != nil {
		return 0, fmt.Errorf("Get error: %v", err)
	}

	var ret struct {
		Error string `json:"error,omitempt"`
		Data  int    `json:"data,omitempt"`
	}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return 0, fmt.Errorf("JSON decode error: %v", err)
	}

	return ret.Data, nil
}

func testAndValidateDeleteReq(t *testing.T, addr string, id int) {
	delID, err := testDeleteReq(addr, id)
	if err != nil {
		t.Errorf("Can't create data: %v", err)
	}
	if delID != id {
		t.Errorf("Expected id of %d, got %d instead.", id, delID)
	}
}

func testBadGet(t *testing.T, addr string, uri string) {
	url := fmt.Sprintf("http://%s%s", addr, uri)
	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected error from %s.  Response: %s", url, string(b))
	}
}

func testBadPost(t *testing.T, addr string, uri string, data interface{}) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)

	url := fmt.Sprintf("http://%s%s", addr, uri)
	resp, err := http.Post(url, "application/json; charset=utf-8", buf)
	if err == nil && resp.StatusCode == http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected error from %s.  Response: %s", url, string(b))
	}
}

func TestNewService(t *testing.T) {
	r := NewRouter()
	s := NewTestService()
	err := r.AddService("test", s)
	if err != nil {
		t.Fatalf("Can't add service: %v", err)
	}

	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatalf("Can't get free port: %v", err)

	}
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	listener, err := net.Listen("tcp", addr)
	srv := &http.Server{
		Handler: r,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go srv.Serve(listener)

	testAndValidateNewReq(t, addr, &TestData{ID: 1, Name: "Test 1"})
	testAndValidateNewReq(t, addr, &TestData{ID: 2, Name: "Test 2"})

	testAndValidateGetReq(t, addr, 1, "Test 1")
	testAndValidateGetReq(t, addr, 2, "Test 2")

	testAndValidatePutReq(t, addr, &TestData{ID: 2, Name: "Test 2 put"})
	testAndValidateGetReq(t, addr, 2, "Test 2 put")

	_, err = testPutReq(addr, 3, &TestData{ID: 3, Name: "Test 3 put"})
	if err == nil {
		t.Errorf("Expected error putting non-existant record")
	}

	testAndValidateDeleteReq(t, addr, 1)
	_, err = testGetReq(addr, 1)
	if err == nil {
		t.Errorf("Expected error after fetching deleted record")
	}

	s.FailNew = true
	_, err = testNewReq(addr, &TestData{ID: 3, Name: "Test 3"})
	if err == nil {
		t.Errorf("Expected error from forced failure of New")
	}
	s.FailNew = false

	_, err = testDeleteReq(addr, 1)
	if err == nil {
		t.Errorf("Expected error from deleting non-existant record")
	}

	testBadGet(t, addr, "/test/get/word")
	testBadGet(t, addr, "/test/get/1000000000000000000000000000000000000000000000001")
	testBadGet(t, addr, "/test/delete/word")
	testBadGet(t, addr, "/test/delete/1000000000000000000000000000000000000000000001")

	testBadPost(t, addr, "/test/put/word", "")
	testBadPost(t, addr, "/test/put/1000000000000000000000000000000000000000000001",
		"")
	testBadPost(t, addr, "/test/put/1", "")

	testBadPost(t, addr, "/test/new", "")
}

func TestEncodeError(t *testing.T) {
	w := NewDummyResponseWriter()
	sendJsonResponse(w, make(chan int))
	if w.StatusCode == http.StatusOK {
		t.Errorf("Expected failure from sendJsonResponse on unencodable data")
	}
}
