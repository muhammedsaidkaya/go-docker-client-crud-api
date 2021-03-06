package controller_layer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/dto"
	"github.com/muhammedsaidkaya/crud-api--container--golang-docker-client/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {

	endpoint := "/containers"
	mockService := mock.NewMockContainerServiceInterface(gomock.NewController(t))

	t.Run("delete method tests", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			id := "ID-123"

			mockService.EXPECT().
				Delete(id).
				Return(nil).
				Times(1)

			router := gin.Default()
			NewContainerController(mockService).Setup(router)
			response := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", endpoint, id), http.NoBody)
			router.ServeHTTP(response, request)

			var jsonMap map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &jsonMap)

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, id, jsonMap["response"].(map[string]interface{})["data"].(string))
			assert.Equal(t, "Container removed.", jsonMap["response"].(map[string]interface{})["message"].(string))
		})

		t.Run("not found", func(t *testing.T) {
			id := "ID-123"

			mockService.EXPECT().
				Delete(id).
				Return(errors.New("not found")).
				Times(1)

			router := gin.Default()
			NewContainerController(mockService).Setup(router)
			response := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", endpoint, id), http.NoBody)
			router.ServeHTTP(response, request)

			var jsonMap map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &jsonMap)

			assert.Equal(t, http.StatusNotFound, response.Code)
			assert.Equal(t, id, jsonMap["response"].(map[string]interface{})["data"].(string))
			assert.Equal(t, "Container not found.", jsonMap["response"].(map[string]interface{})["message"].(string))
		})

	})

	t.Run("getByID method tests", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			id := "ID-123"
			expected := dto.ContainerMetadata{ID: id, Names: []string{"uzumlukek-test"}, Image: "uzumlukek/docker-client"}

			mockService.EXPECT().
				GetByID(id).
				Return(types.Container{ID: id, Names: []string{"uzumlukek-test"}, Image: "uzumlukek/docker-client"}, nil).
				Times(1)

			router := gin.Default()
			NewContainerController(mockService).Setup(router)
			response := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", endpoint, id), http.NoBody)
			router.ServeHTTP(response, request)

			var jsonMap map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &jsonMap)

			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, expected.ID, jsonMap["response"].(map[string]interface{})["data"].(map[string]interface{})["Id"].(string))
		})

		t.Run("not found", func(t *testing.T) {
			id := "ID-123"
			mockService.EXPECT().
				GetByID(id).
				Return(types.Container{}, errors.New("not found")).
				Times(1)

			router := gin.Default()
			NewContainerController(mockService).Setup(router)
			response := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", endpoint, id), http.NoBody)
			router.ServeHTTP(response, request)

			var jsonMap map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &jsonMap)

			assert.Equal(t, http.StatusNotFound, response.Code)
			assert.Nil(t, jsonMap["response"].(map[string]interface{})["data"])
		})
	})
}
