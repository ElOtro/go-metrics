package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElOtro/go-metrics/internal/handlers/mocks"
	"github.com/ElOtro/go-metrics/internal/repo/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestHandlers_GetAllMetricsHandler(t *testing.T) {
	type fields struct {
		r    *chi.Mux
		repo *mocks.Repo
	}
	// определяем структуру теста
	type want struct {
		statusCode      int
		wantCallService bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		// определяем все тесты
		{
			name: "Test 1",
			fields: fields{
				r:    chi.NewRouter(),
				repo: &mocks.Repo{},
			},
			want: want{
				statusCode:      http.StatusOK,
				wantCallService: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//  если в процессе теста вызываается сервис мокаем
			if tt.want.wantCallService {
				tt.fields.repo.On("GetAll").Return(make(map[string]float64), make(map[string]int64))
			}

			h := &Handlers{
				repo: tt.fields.repo,
			}

			request := httptest.NewRequest(http.MethodGet, "/", nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()

			// определяем хендлер
			hh := http.HandlerFunc(h.GetAllMetricsHandler)
			hh.ServeHTTP(w, request)
			res := w.Result()
			// получаем и проверяем тело запроса
			defer res.Body.Close()

			// проверяем код ответа
			if res.StatusCode != tt.want.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.want.statusCode, w.Code)
			}

			{
				assert.Equal(t, tt.want.statusCode, res.StatusCode)
			}

			// проверяем что метод у мокового сервиса вызывался
			if tt.want.wantCallService {
				tt.fields.repo.AssertCalled(t, "GetAll")
				tt.fields.repo.AssertNumberOfCalls(t, "GetAll", 1)
			}
		})
	}
}

func TestHandlers_GetMetricHandler(t *testing.T) {
	type fields struct {
		repo Repo
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				repo: tt.fields.repo,
			}
			h.GetMetricHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestHandlers_CreateMetricsJSONHandler(t *testing.T) {
	type fields struct {
		r    *chi.Mux
		repo *mocks.Repo
	}
	// определяем структуру теста
	type want struct {
		statusCode      int
		wantCallService bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{

		// определяем все тесты
		{
			name: "Create new metric",
			fields: fields{
				r:    chi.NewRouter(),
				repo: &mocks.Repo{},
			},
			want: want{
				statusCode:      http.StatusOK,
				wantCallService: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var value float64 = 3824
			metric := &storage.Metrics{ID: "BuckHashSys", MType: "gauge", Value: &value}

			if tt.want.wantCallService {
				tt.fields.repo.On("SetMetrics", metric).Return(nil)
			}

			h := &Handlers{
				repo: tt.fields.repo,
			}

			js := []byte(`{"id":"BuckHashSys","type":"gauge","value":3824}`)
			body := bytes.NewReader(js)

			request := httptest.NewRequest(http.MethodPost, "/update/", body)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			hh := http.HandlerFunc(h.CreateMetricsJSONHandler)
			hh.ServeHTTP(w, request)
			res := w.Result()
			// получаем и проверяем тело запроса
			defer res.Body.Close()

			// проверяем код ответа
			if res.StatusCode != tt.want.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.want.statusCode, w.Code)
			}

			{
				assert.Equal(t, tt.want.statusCode, res.StatusCode)
			}

			if tt.want.wantCallService {
				tt.fields.repo.AssertCalled(t, "SetMetrics", metric)
				tt.fields.repo.AssertNumberOfCalls(t, "SetMetrics", 1)
			}

		})
	}
}

func TestHandlers_GetMetricsJSONHandler(t *testing.T) {
	type fields struct {
		repo *mocks.Repo
	}
	// определяем структуру теста
	type want struct {
		statusCode int
		value      float64
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "Get JSON metric",
			fields: fields{
				repo: &mocks.Repo{},
			},
			want: want{
				statusCode: http.StatusOK,
				value:      3824,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var value float64 = 3824
			metric := &storage.Metrics{ID: "BuckHashSys", MType: "gauge", Value: &value}

			tt.fields.repo.On("GetMetricsByID", metric.ID, metric.MType).Return(metric, nil)

			h := &Handlers{
				repo: tt.fields.repo,
			}

			js := []byte(`{"id":"BuckHashSys","type":"gauge"}`)
			body := bytes.NewReader(js)

			request := httptest.NewRequest(http.MethodPost, "/value/", body)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			hh := http.HandlerFunc(h.GetMetricsJSONHandler)
			hh.ServeHTTP(w, request)
			res := w.Result()
			// получаем и проверяем тело запроса
			defer res.Body.Close()

			fmt.Println(res)

			// проверяем код ответа
			if res.StatusCode != tt.want.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.want.statusCode, w.Code)
			}

			{
				assert.Equal(t, tt.want.statusCode, res.StatusCode)
				assert.Equal(t, tt.want.value, value)
			}

			tt.fields.repo.AssertCalled(t, "GetMetricsByID", metric.ID, metric.MType)
			tt.fields.repo.AssertNumberOfCalls(t, "GetMetricsByID", 1)

		})
	}
}
