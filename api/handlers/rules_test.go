package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/odpf/siren/api/handlers"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/logger"
	"github.com/odpf/siren/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func getPanicLogger() *zap.Logger {
	panicLogger, _ := logger.New(&domain.LogConfig{Level: "panic"})
	return panicLogger
}

func TestRules_UpsertRules(t *testing.T) {
	t.Run("should return 200 OK on success", func(t *testing.T) {
		mockedRulesService := &mocks.RuleService{}
		dummyRule := &domain.Rule{
			Namespace: "foo",
			Entity:    "gojek", GroupName: "test-group", Template: "test-tmpl", Status: "enabled",
			Variables: []domain.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			},
			},
		}
		payload := []byte(`{"namespace":"foo","group_name":"test-group","entity":"gojek","template":"test-tmpl","status":"enabled", "variables": [{"name": "test-name", "value":"test-value", "description": "test-description", "type": "test-type" }]}`)

		mockedRulesService.On("Upsert", dummyRule).Return(dummyRule, nil).Once()
		r, err := http.NewRequest(http.MethodPut, "/rules", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.UpsertRule(mockedRulesService, getPanicLogger())
		expectedStatusCode := http.StatusOK
		response, _ := json.Marshal(dummyRule)
		expectedStringBody := string(response) + "\n"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedRulesService.AssertCalled(t, "Upsert", dummyRule)
	})

	t.Run("should return 400 Bad Request on failure", func(t *testing.T) {
		mockedRulesService := &mocks.RuleService{}
		dummyRule := &domain.Rule{
			Namespace: "foo",
			Entity:    "gojek", GroupName: "test-group", Template: "test-tmpl", Status: "enabled",
			Variables: []domain.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			},
			},
		}
		payload := []byte(`bad input`)

		mockedRulesService.On("Upsert", dummyRule).Return(dummyRule, nil).Once()
		r, err := http.NewRequest(http.MethodPut, "/rules", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.UpsertRule(mockedRulesService, getPanicLogger())
		expectedStatusCode := http.StatusBadRequest
		expectedStringBody := "{\"code\":400,\"message\":\"invalid character 'b' looking for beginning of value\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedRulesService.AssertNotCalled(t, "Upsert", dummyRule)
	})

	t.Run("should return 400 Bad Request if request body validation fails", func(t *testing.T) {
		mockedRulesService := &mocks.RuleService{}
		payload := []byte(`{"namespace":"","group_name":"test-group","entity":"gojek","template":"test-tmpl","status":"enabled", "variables": [{"name": "test-name", "value":"test-value", "description": "test-description", "type": "test-type" }]}`)

		r, err := http.NewRequest(http.MethodPut, "/rules", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.UpsertRule(mockedRulesService, getPanicLogger())
		expectedStatusCode := http.StatusBadRequest
		expectedStringBody := "{\"code\":400,\"message\":\"Key: 'Rule.Namespace' Error:Field validation for 'Namespace' failed on the 'required' tag\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedRulesService.AssertNotCalled(t, "Upsert")
	})

	t.Run("should return 400 Bad Request if template not found", func(t *testing.T) {
		mockedRulesService := &mocks.RuleService{}
		dummyRule := &domain.Rule{
			Namespace: "foo",
			Entity:    "gojek", GroupName: "test-group", Template: "test-tmpl", Status: "enabled",
			Variables: []domain.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			},
			},
		}
		payload := []byte(`{"namespace":"foo","group_name":"test-group","entity":"gojek","template":"test-tmpl","status":"enabled", "variables": [{"name": "test-name", "value":"test-value", "description": "test-description", "type": "test-type" }]}`)

		mockedRulesService.On("Upsert", dummyRule).Return(nil, errors.New("template not found")).Once()
		r, err := http.NewRequest(http.MethodPut, "/rules", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.UpsertRule(mockedRulesService, getPanicLogger())
		expectedStatusCode := http.StatusBadRequest
		expectedStringBody := "{\"code\":400,\"message\":\"template not found\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})

	t.Run("should return 500 Internal Server Error on service failure", func(t *testing.T) {
		mockedRulesService := &mocks.RuleService{}
		dummyRule := &domain.Rule{
			Namespace: "foo",
			Entity:    "gojek", GroupName: "test-group", Template: "test-tmpl", Status: "enabled",
			Variables: []domain.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			},
			},
		}
		payload := []byte(`{"namespace":"foo","group_name":"test-group","entity":"gojek","template":"test-tmpl","status":"enabled", "variables": [{"name": "test-name", "value":"test-value", "description": "test-description", "type": "test-type" }]}`)

		mockedRulesService.On("Upsert", dummyRule).Return(nil, errors.New("random error")).Once()
		r, err := http.NewRequest(http.MethodPut, "/rules", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.UpsertRule(mockedRulesService, getPanicLogger())
		expectedStatusCode := http.StatusInternalServerError
		expectedStringBody := "{\"code\":500,\"message\":\"Internal server error\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedRulesService.AssertCalled(t, "Upsert", dummyRule)
	})
}

func TestRules_GetRules(t *testing.T) {
	t.Run("should return 200 OK on success", func(t *testing.T) {
		mockedRulesService := &mocks.RuleService{}
		dummyRules := []domain.Rule{{
			Namespace: "foo",
			Entity:    "gojek", GroupName: "test-group", Template: "test-tmpl", Status: "enabled",
			Variables: []domain.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			},
			},
		}}

		mockedRulesService.On("Get", "foo", "gojek", "bar", "enabled", "tmpl").Return(dummyRules, nil).Once()
		r, err := http.NewRequest(http.MethodGet, "/rules", nil)
		q := r.URL.Query()
		q.Add("namespace", "foo")
		q.Add("entity", "gojek")
		q.Add("group_name", "bar")
		q.Add("status", "enabled")
		q.Add("template", "tmpl")
		r.URL.RawQuery = q.Encode()
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.GetRules(mockedRulesService, getPanicLogger())
		expectedStatusCode := http.StatusOK
		response, _ := json.Marshal(dummyRules)
		expectedStringBody := string(response) + "\n"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedRulesService.AssertCalled(t, "Get", "foo", "gojek", "bar", "enabled", "tmpl")
	})

	t.Run("should return 500 Internal Server Error on service failure", func(t *testing.T) {
		mockedRulesService := &mocks.RuleService{}

		mockedRulesService.On("Get", "foo", "gojek", "bar", "enabled", "tmpl").Return(nil, errors.New("random error")).Once()
		r, err := http.NewRequest(http.MethodGet, "/rules", nil)
		q := r.URL.Query()
		q.Add("namespace", "foo")
		q.Add("entity", "gojek")
		q.Add("group_name", "bar")
		q.Add("status", "enabled")
		q.Add("template", "tmpl")
		r.URL.RawQuery = q.Encode()
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.GetRules(mockedRulesService, getPanicLogger())
		expectedStatusCode := http.StatusInternalServerError
		expectedStringBody := "{\"code\":500,\"message\":\"Internal server error\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedRulesService.AssertCalled(t, "Get", "foo", "gojek", "bar", "enabled", "tmpl")
	})
}
