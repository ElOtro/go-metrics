package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ElOtro/go-metrics/cmd/handlers/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestHandlers_GetAllMetricsHandler(t *testing.T) {
	type fields struct {
		r    *chi.Mux
		repo *mocks.Repo
	}
	type args struct {
		value string
		// w http.ResponseWriter
		// r *http.Request
	}
	// определяем структуру теста
	type want struct {
		statusCode      int
		value           string
		wantCallService bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		// определяем все тесты
		{
			name: "Test 1",
			fields: fields{
				r:    chi.NewRouter(),
				repo: &mocks.Repo{},
			},
			args: args{
				value: "/",
			},
			want: want{
				statusCode:      http.StatusOK,
				value:           "4",
				wantCallService: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//  если в процессе теста вызываается сервсис мокаем
			if tt.want.wantCallService {
				val, err := strconv.Atoi(tt.args.value)
				assert.NoError(t, err)
				res, err := strconv.Atoi(tt.want.value)
				assert.NoError(t, err)
				tt.fields.repo.On("Repo", val).Return(res)
			}

			h := &Handlers{
				repo: tt.fields.repo,
			}

			url := fmt.Sprintf("http://%s/%s", httptest.DefaultRemoteAddr, tt.args.value)
			request := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			hh := http.HandlerFunc(h.GetAllMetricsHandler)
			hh.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.want.statusCode {
				t.Errorf("Status code mismatched got %v, want %v", res.StatusCode, tt.want.statusCode)
			}

			{
				result, err := ioutil.ReadAll(res.Body)
				assert.NoError(t, err)
				assert.Equal(t, tt.want.value, string(result))
			}
			// проверяем что метод у мокового сервиса вызывался
			if tt.want.wantCallService {
				val, err := strconv.Atoi(tt.args.value)
				assert.NoError(t, err)
				tt.fields.repo.AssertCalled(t, "Repo", val)
				tt.fields.repo.AssertNumberOfCalls(t, "Repo", 1)
			}
		})
	}
}
