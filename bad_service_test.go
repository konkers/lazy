package lazy

import (
	"context"
	"net/url"
	"testing"
)

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

type TestServiceBadGetIn0 struct {
	TestService
}

func (s *TestServiceBadGetIn0) Get(ctx string, id int) {
}

type TestServiceBadGetIn1 struct {
	TestService
}

func (s *TestServiceBadGetIn1) Get(ctx context.Context, id string) {
}

type TestServiceBadGetOutNum struct {
	TestService
}

func (s *TestServiceBadGetOutNum) Get(ctx context.Context, id int) {
}

type TestServiceBadGetOut0 struct {
	TestService
}

func (s *TestServiceBadGetOut0) Get(ctx context.Context, id int) (int, error) {
	return 0, nil
}

type TestServiceBadGetOut0Private struct {
	TestService
}

func (s *TestServiceBadGetOut0Private) Get(ctx context.Context, id int) (*privateData, error) {
	return nil, nil
}

type TestServiceBadGetOut1 struct {
	TestService
}

func (s *TestServiceBadGetOut1) Get(ctx context.Context, id int) (*TestData, int) {
	return nil, 0
}

//
// Put
//

type TestServiceNoPut struct {
}

func (s *TestServiceNoPut) Get(ctx context.Context, id int) (*TestData, error) {
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

func (s *TestServiceBadPutIn0) Put(ctx string, id int, data *TestData) {
}

type TestServiceBadPutIn1 struct {
	TestService
}

func (s *TestServiceBadPutIn1) Put(ctx context.Context, id string, data *TestData) {
}

type TestServiceBadPutIn2 struct {
	TestService
}

func (s *TestServiceBadPutIn2) Put(ctx context.Context, id int, data int) {
}

type TestServiceBadPutOutNum struct {
	TestService
}

func (s *TestServiceBadPutOutNum) Put(ctx context.Context, id int, data *TestData) {
}

type TestServiceBadPutOut0 struct {
	TestService
}

func (s *TestServiceBadPutOut0) Put(ctx context.Context, id int, data *TestData) int {
	return 0
}

//
// New
//

type TestServiceNoNew struct {
}

func (s *TestServiceNoNew) Get(ctx context.Context, id int) (*TestData, error) {
	return nil, nil
}

func (s *TestServiceNoNew) Put(ctx context.Context, id int, data *TestData) error {
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

func (s *TestServiceBadNewIn0) New(ctx int, data *TestData) {
}

type TestServiceBadNewIn1 struct {
	TestService
}

func (s *TestServiceBadNewIn1) New(ctx context.Context, data int) {
}

type TestServiceBadNewOutNum struct {
	TestService
}

func (s *TestServiceBadNewOutNum) New(ctx context.Context, data *TestData) {
}

type TestServiceBadNewOut0 struct {
	TestService
}

func (s *TestServiceBadNewOut0) New(ctx context.Context, data *TestData) (string, error) {
	return "", nil
}

type TestServiceBadNewOut1 struct {
	TestService
}

func (s *TestServiceBadNewOut1) New(ctx context.Context, data *TestData) (int, string) {
	return 0, ""
}

//
// Delete
//

type TestServiceNoDelete struct {
}

func (s *TestServiceNoDelete) Get(ctx context.Context, id int) (*TestData, error) {
	return nil, nil
}

func (s *TestServiceNoDelete) Put(ctx context.Context, id int, data *TestData) error {
	return nil
}

func (s *TestServiceNoDelete) New(ctx context.Context, data *TestData) (int, error) {
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

func (s *TestServiceBadDeleteIn0) Delete(ctx int, id int) {
}

type TestServiceBadDeleteIn1 struct {
	TestService
}

func (s *TestServiceBadDeleteIn1) Delete(ctx context.Context, id string) {
}

type TestServiceBadDeleteOutNum struct {
	TestService
}

func (s *TestServiceBadDeleteOutNum) Delete(ctx context.Context, id int) {
}

type TestServiceBadDeleteOut0 struct {
	TestService
}

func (s *TestServiceBadDeleteOut0) Delete(ctx context.Context, id int) string {
	return ""
}

//
// Query
//

type TestServiceNoQuery struct {
}

func (s *TestServiceNoQuery) Get(ctx context.Context, id int) (*TestData, error) {
	return nil, nil
}

func (s *TestServiceNoQuery) Put(ctx context.Context, id int, data *TestData) error {
	return nil
}

func (s *TestServiceNoQuery) New(ctx context.Context, data *TestData) (int, error) {
	return 0, nil
}
func (s *TestServiceNoQuery) Delete(ctx context.Context, id int) error {
	return nil
}

type TestServiceBadQueryInNum struct {
	TestService
}

func (s *TestServiceBadQueryInNum) Query() {
}

type TestServiceBadQueryIn0 struct {
	TestService
}

func (s *TestServiceBadQueryIn0) Query(ctx string, args url.Values) ([]*TestData, error) {
	return nil, nil
}

type TestServiceBadQueryIn1 struct {
	TestService
}

func (s *TestServiceBadQueryIn1) Query(ctx context.Context, args string) ([]*TestData, error) {
	return nil, nil
}

type TestServiceBadQueryOutNum struct {
	TestService
}

func (s *TestServiceBadQueryOutNum) Query(ctx context.Context, args url.Values) {
}

type TestServiceBadQueryOut0 struct {
	TestService
}

func (s *TestServiceBadQueryOut0) Query(ctx context.Context, args url.Values) (string, error) {
	return "", nil
}

type TestServiceBadQueryOut1 struct {
	TestService
}

func (s *TestServiceBadQueryOut1) Query(ctx context.Context, args url.Values) ([]*TestData, string) {
	return nil, ""
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
		{&TestServiceBadGetIn0{}, "Expected error from bad Get argument 0 type"},
		{&TestServiceBadGetIn1{}, "Expected error from bad Get argument 1 type"},
		{&TestServiceBadGetOutNum{}, "Expected error from bad Get out arg count"},
		{&TestServiceBadGetOut0{}, "Expected error from bad Get out arg 0"},
		{&TestServiceBadGetOut0Private{}, "Expected error from private Get out arg 0"},
		{&TestServiceBadGetOut1{}, "Expected error from bad Get out arg 1"},

		// Put
		{&TestServiceNoPut{}, "Expected error from service without Put"},
		{&TestServiceBadPutInNum{}, "Expected error from bad Put in arg count"},
		{&TestServiceBadPutIn0{}, "Expected error from bad Put argument 0 type"},
		{&TestServiceBadPutIn1{}, "Expected error from bad Put argument 1 type"},
		{&TestServiceBadPutIn2{}, "Expected error from bad Put argument 2 type"},
		{&TestServiceBadPutOutNum{}, "Expected error from bad Put out arg count"},
		{&TestServiceBadPutOut0{}, "Expected error from bad Put out arg 0"},

		// New
		{&TestServiceNoNew{}, "Expected error from service without New"},
		{&TestServiceBadNewInNum{}, "Expected error from bad New in arg count"},
		{&TestServiceBadNewIn0{}, "Expected error from bad New argument 0 type"},
		{&TestServiceBadNewIn1{}, "Expected error from bad New argument 1 type"},
		{&TestServiceBadNewOutNum{}, "Expected error from bad New out arg count"},
		{&TestServiceBadNewOut0{}, "Expected error from bad New out arg 0"},
		{&TestServiceBadNewOut1{}, "Expected error from bad New out arg 1"},

		// Delete
		{&TestServiceNoDelete{}, "Expected error from service without Delete"},
		{&TestServiceBadDeleteInNum{}, "Expected error from bad Delete in arg count"},
		{&TestServiceBadDeleteIn0{}, "Expected error from bad Delete argument 0 type"},
		{&TestServiceBadDeleteIn1{}, "Expected error from bad Delete argument 1 type"},
		{&TestServiceBadDeleteOutNum{}, "Expected error from bad Query out arg count"},
		{&TestServiceBadDeleteOut0{}, "Expected error from bad Delete out arg 0"},

		// Query
		{&TestServiceNoQuery{}, "Expected error from service without Query"},
		{&TestServiceBadQueryInNum{}, "Expected error from bad Query in arg count"},
		{&TestServiceBadQueryIn0{}, "Expected error from bad Query argument 0 type"},
		{&TestServiceBadQueryIn1{}, "Expected error from bad Query argument 1 type"},
		{&TestServiceBadQueryOutNum{}, "Expected error from bad Query out arg count"},
		{&TestServiceBadQueryOut0{}, "Expected error from bad Query out arg 0"},
		{&TestServiceBadQueryOut1{}, "Expected error from bad Query out arg 1"},
	}

	for _, test := range tests {
		err := r.AddService("test", test.service)
		if err == nil {
			t.Errorf(test.errorMsg)
		}
	}
}
