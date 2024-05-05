package shorten

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/mocks"
)

const URL = "www.yandex.ru"
const KEY = "key"

func TestShortenService(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	storage := mocks.NewMockRepository(mockCtrl)
	storage.EXPECT().FindByKey(KEY).Return(URL, true)
	service := NewURLService(storage, &configuration.Config{})
	val, ok := service.FindByKey(KEY)
	assert.Equal(t, val, URL)
	assert.Equal(t, ok, true)
}
