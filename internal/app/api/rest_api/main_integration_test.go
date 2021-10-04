package rest_api_test

import (
	"github.com/VSKrivoshein/FBS-test/internal/app/api/rest_api"
	"github.com/VSKrivoshein/FBS-test/internal/app/configs"
	"net/http/httptest"
	"os"
	"testing"
)

var TestSrv *httptest.Server

func TestMain(m *testing.M) {
	service := configs.InitService()
	handler := rest_api.NewHandler(service).InitRoutes()
	TestSrv = httptest.NewServer(handler)
	defer TestSrv.Close()

	exitVal := m.Run()
	os.Exit(exitVal)
}
