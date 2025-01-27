package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfig_AddDefaultData(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	testApp.Session.Put(ctx, "flash", "flash")
	testApp.Session.Put(ctx, "warning", "warning")
	testApp.Session.Put(ctx, "error", "error")

	td := testApp.AddDefaultData(&TemplateData{}, req)

	if td.Flash != "flash" {
		t.Error("Failed to get flash data!")
	}
	if td.Warning != "warning" {
		t.Error("Failed to get warning data!")
	}
	if td.Error != "error" {
		t.Error("Failed to get error data!")
	}
}

func TestConfig_IsAuthenticated(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	auth := testApp.IsAuthenticate(req)
	if auth {
		t.Error("Return TRUE for authenticated, when it should be FALSE!")
	}

	testApp.Session.Put(ctx, "userID", 1)
	auth = testApp.IsAuthenticate(req)
	if !auth {
		t.Error("Return FALSE for authenticated, when it should be TRUE!")
	}
}

func TestConfig_render(t *testing.T) {
	pathToTemplates = "./templates"

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	testApp.render(rr, req, "home.page.gohtml", &TemplateData{})

	if rr.Code != 200 {
		t.Error("Failed to render page!")
	}
}
