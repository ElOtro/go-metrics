package handlers

import (
	"net/http"
	"testing"
)

// func TestHandlers_GetAllMetricsHandler(t *testing.T) {
// 	type fields struct {
// 		r    *chi.Mux
// 		repo *mocks.Repo
// 	}
// 	// определяем структуру теста
// 	type want struct {
// 		statusCode      int
// 		wantCallService bool
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   want
// 	}{
// 		// определяем все тесты
// 		{
// 			name: "Test 1",
// 			fields: fields{
// 				r:    chi.NewRouter(),
// 				repo: &mocks.Repo{},
// 			},
// 			want: want{
// 				statusCode:      http.StatusOK,
// 				wantCallService: true,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			//  если в процессе теста вызываается сервис мокаем
// 			if tt.want.wantCallService {
// 				tt.fields.repo.On("GetAll")
// 			}

// 			h := &Handlers{
// 				repo: tt.fields.repo,
// 			}

// 			request := httptest.NewRequest(http.MethodGet, "/", nil)

// 			// создаём новый Recorder
// 			w := httptest.NewRecorder()

// 			// определяем хендлер
// 			hh := http.HandlerFunc(h.GetAllMetricsHandler)
// 			hh.ServeHTTP(w, request)
// 			res := w.Result()
// 			// получаем и проверяем тело запроса
// 			defer res.Body.Close()

// 			// проверяем код ответа
// 			if res.StatusCode != tt.want.statusCode {
// 				t.Errorf("Expected status code %d, got %d", tt.want.statusCode, w.Code)
// 			}

// 			{
// 				assert.Equal(t, tt.want.statusCode, res.StatusCode)
// 			}

// 			// проверяем что метод у мокового сервиса вызывался
// 			if tt.want.wantCallService {
// 				tt.fields.repo.AssertCalled(t, "GetAll")
// 				tt.fields.repo.AssertNumberOfCalls(t, "GetAll", 1)
// 			}
// 		})
// 	}
// }

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
