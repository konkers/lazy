package lazy

import "testing"

type privateData struct{}

//
// Get
//

type TestServiceNoGet struct {
}

type TestServiceBadGetInNum struct {
	TestService
}

func (s *TestServiceBadGetInNum) Get() {
}

type TestServiceBadGetIn1 struct {
	TestService
}

func (s *TestServiceBadGetIn1) Get(id string) {
}

type TestServiceBadGetOutNum struct {
	TestService
}

func (s *TestServiceBadGetOutNum) Get(id int) {
}

type TestServiceBadGetOut0 struct {
	TestService
}

func (s *TestServiceBadGetOut0) Get(id int) (int, error) {
	return 0, nil
}

type TestServiceBadGetOut0Private struct {
	TestService
}

func (s *TestServiceBadGetOut0Private) Get(id int) (*privateData, error) {
	return nil, nil
}

type TestServiceBadGetOut1 struct {
	TestService
}

func (s *TestServiceBadGetOut1) Get(id int) (*TestData, int) {
	return nil, 0
}

//
// Put
//

type TestServiceNoPut struct {
}

func (s *TestServiceNoPut) Get(id int) (*TestData, error) {
	return nil, nil
}

type TestServiceBadPutInNum struct {
	TestService
}

func (s *TestServiceBadPutInNum) Put() {
}

type TestServiceBadPutIn0 struct {
	TestService
}

func (s *TestServiceBadPutIn0) Put(id string, data *TestData) {
}

type TestServiceBadPutIn1 struct {
	TestService
}

func (s *TestServiceBadPutIn1) Put(id int, data int) {
}

type TestServiceBadPutOutNum struct {
	TestService
}

func (s *TestServiceBadPutOutNum) Put(id int, data *TestData) {
}

type TestServiceBadPutOut0 struct {
	TestService
}

func (s *TestServiceBadPutOut0) Put(id int, data *TestData) int {
	return 0
}

//
// New
//

type TestServiceNoNew struct {
}

func (s *TestServiceNoNew) Get(id int) (*TestData, error) {
	return nil, nil
}

func (s *TestServiceNoNew) Put(id int, data *TestData) error {
	return nil
}

type TestServiceBadNewInNum struct {
	TestService
}

func (s *TestServiceBadNewInNum) New() {
}

type TestServiceBadNewIn0 struct {
	TestService
}

func (s *TestServiceBadNewIn0) New(data int) {
}

type TestServiceBadNewOutNum struct {
	TestService
}

func (s *TestServiceBadNewOutNum) New(data *TestData) {
}

type TestServiceBadNewOut0 struct {
	TestService
}

func (s *TestServiceBadNewOut0) New(data *TestData) (string, error) {
	return "", nil
}

type TestServiceBadNewOut1 struct {
	TestService
}

func (s *TestServiceBadNewOut1) New(data *TestData) (int, string) {
	return 0, ""
}

//
// New
//

type TestServiceNoDelete struct {
}

func (s *TestServiceNoDelete) Get(id int) (*TestData, error) {
	return nil, nil
}

func (s *TestServiceNoDelete) Put(id int, data *TestData) error {
	return nil
}

func (s *TestServiceNoDelete) New(data *TestData) (int, error) {
	return 0, nil
}

type TestServiceBadDeleteInNum struct {
	TestService
}

func (s *TestServiceBadDeleteInNum) Delete() {
}

type TestServiceBadDeleteIn0 struct {
	TestService
}

func (s *TestServiceBadDeleteIn0) Delete(id string) {
}

type TestServiceBadDeleteOutNum struct {
	TestService
}

func (s *TestServiceBadDeleteOutNum) Delete(id int) {
}

type TestServiceBadDeleteOut0 struct {
	TestService
}

func (s *TestServiceBadDeleteOut0) Delete(id int) string {
	return ""
}

func TestBadService(t *testing.T) {
	r := NewRouter()

	type testdata struct {
		service  interface{}
		errorMsg string
	}
	tests := []testdata{
		// Get
		{&TestServiceNoGet{}, "Expected error from service without Get"},
		{&TestServiceBadGetInNum{}, "Expected error from bad Get in arg count"},
		{&TestServiceBadGetIn1{}, "Expected error from bad Get argument type"},
		{&TestServiceBadGetOutNum{}, "Expected error from bad Get out arg count"},
		{&TestServiceBadGetOut0{}, "Expected error from bad Get out arg 0"},
		{&TestServiceBadGetOut0Private{}, "Expected error from private Get out arg 0"},
		{&TestServiceBadGetOut1{}, "Expected error from bad Get out arg 1"},

		// Put
		{&TestServiceNoPut{}, "Expected error from service without Put"},
		{&TestServiceBadPutInNum{}, "Expected error from bad Put in arg count"},
		{&TestServiceBadPutIn0{}, "Expected error from bad Put argument 0 type"},
		{&TestServiceBadPutIn1{}, "Expected error from bad Put argument 1 type"},
		{&TestServiceBadPutOutNum{}, "Expected error from bad Put out arg count"},
		{&TestServiceBadPutOut0{}, "Expected error from bad Put out arg 0"},

		// New
		{&TestServiceNoNew{}, "Expected error from service without New"},
		{&TestServiceBadNewInNum{}, "Expected error from bad New in arg count"},
		{&TestServiceBadNewIn0{}, "Expected error from bad New argument 0 type"},
		{&TestServiceBadNewOutNum{}, "Expected error from bad New out arg count"},
		{&TestServiceBadNewOut0{}, "Expected error from bad New out arg 0"},
		{&TestServiceBadNewOut1{}, "Expected error from bad New out arg 1"},

		// Delete
		{&TestServiceNoDelete{}, "Expected error from service without Delete"},
		{&TestServiceBadDeleteInNum{}, "Expected error from bad Delete in arg count"},
		{&TestServiceBadDeleteIn0{}, "Expected error from bad Delete argument 0 type"},
		{&TestServiceBadDeleteOutNum{}, "Expected error from bad New out arg count"},
		{&TestServiceBadDeleteOut0{}, "Expected error from bad Delete out arg 0"},
	}

	for _, test := range tests {
		err := r.AddService("test", test.service)
		if err == nil {
			t.Errorf(test.errorMsg)
		}
	}
}
