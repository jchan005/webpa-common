package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRedirectHandler(t *testing.T) {
	var (
		assert       = assert.New(t)
		require      = require.New(t)
		expectedKey  = []byte("here is a lovely little key, full of vim and vigour")
		expectedNode = "http://nogohere.instead.com"

		mockAccessor = new(mockAccessor)
		response     = httptest.NewRecorder()
		request      = httptest.NewRequest("GET", "http://foobar.com/test", nil)
		keyFunc      = func(actualRequest *http.Request) ([]byte, error) {
			assert.Equal(request, actualRequest)
			return expectedKey, nil
		}
	)

	mockAccessor.On("Get", expectedKey).Return(expectedNode, nil).Once()

	handler := NewRedirectHandler(mockAccessor, http.StatusTemporaryRedirect, keyFunc, nil)
	handler.ServeHTTP(response, request)
	assert.Equal(http.StatusTemporaryRedirect, response.Code)

	actualLocation, err := response.Result().Location()
	require.NotNil(actualLocation)
	assert.Equal("http://nogohere.instead.com/test", actualLocation.String())
	assert.Nil(err)

	mockAccessor.AssertExpectations(t)
}

func TestNewRedirectHandlerBadKey(t *testing.T) {
	assert := assert.New(t)
	expectedKeyError := errors.New("Expected error from key function")

	mockAccessor := new(mockAccessor)
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://foobar.com/test", nil)
	keyFunc := func(actualRequest *http.Request) ([]byte, error) {
		assert.Equal(request, actualRequest)
		return []byte{}, expectedKeyError
	}

	handler := NewRedirectHandler(mockAccessor, http.StatusTemporaryRedirect, keyFunc, nil)
	handler.ServeHTTP(response, request)
	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Contains(response.Body.String(), expectedKeyError.Error())

	actualLocation, err := response.Result().Location()
	assert.Nil(actualLocation)
	assert.NotNil(err)

	mockAccessor.AssertExpectations(t)
}

func TestNewRedirectHandlerNoNode(t *testing.T) {
	assert := assert.New(t)
	expectedKey := []byte("this little key went to market ...")
	expectedAccessorError := errors.New("Expected error from the Accessor")

	mockAccessor := new(mockAccessor)
	mockAccessor.On("Get", expectedKey).Return("", expectedAccessorError).Once()
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://foobar.com/test", nil)
	keyFunc := func(actualRequest *http.Request) ([]byte, error) {
		assert.Equal(request, actualRequest)
		return expectedKey, nil
	}

	handler := NewRedirectHandler(mockAccessor, http.StatusTemporaryRedirect, keyFunc, nil)
	handler.ServeHTTP(response, request)
	assert.Equal(http.StatusInternalServerError, response.Code)
	assert.Contains(response.Body.String(), expectedAccessorError.Error())

	actualLocation, err := response.Result().Location()
	assert.Nil(actualLocation)
	assert.NotNil(err)

	mockAccessor.AssertExpectations(t)
}
