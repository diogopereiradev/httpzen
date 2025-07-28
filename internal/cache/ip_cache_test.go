package ip_cache_module

import (
	"os"
	"path/filepath"
	"testing"

	app_path_util "github.com/diogopereiradev/httpzen/internal/utils/app_path"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupTestCache() func() {
	CacheFileName = "ip_cache_test"
	cacheViper = newCacheViper()
	configPath := app_path_util.GetConfigPath()
	_ = os.MkdirAll(configPath, 0755)
	cacheFile := filepath.Join(configPath, "ip_cache_test.json")
	_ = os.Remove(cacheFile)
	return func() {
		_ = os.Remove(cacheFile)
	}
}

func TestGetIpInfoFromCache_SetIpInfoToCache(t *testing.T) {
	teardown := setupTestCache()
	defer teardown()

	ip := "127.0.0.1"
	info := map[string]any{"country": "BR", "city": "São Paulo"}

	SetIpInfoToCache(ip, info)

	result, ok := GetIpInfoFromCache("127.0.0.1")
	assert.True(t, ok)
	assert.Equal(t, info, result)
}

func TestGetIpInfoFromCache_NotFound(t *testing.T) {
	teardown := setupTestCache()
	defer teardown()

	result, ok := GetIpInfoFromCache("8.8.8.8")
	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestSetIpInfoToCache_FilePersistence(t *testing.T) {
	teardown := setupTestCache()
	defer teardown()

	ip := "10.0.0.1"
	info := map[string]any{"country": "US", "city": "New York"}
	SetIpInfoToCache(ip, info)

	configPath := app_path_util.GetConfigPath()
	cacheFile := filepath.Join(configPath, "ip_cache_test.json")
	v := viper.New()
	v.SetConfigFile(cacheFile)
	err := v.ReadInConfig()

	assert.NoError(t, err, "Error reading cache file")
	_ = v.MergeInConfig()

	stored := v.GetStringMap("10_0_0_1")
	assert.Equal(t, info, stored)
}

func TestClearCache(t *testing.T) {
	teardown := setupTestCache()
	defer teardown()

	ip := "10.0.0.1"
	info := map[string]any{"country": "US", "city": "New York"}
	SetIpInfoToCache(ip, info)

	ClearCache()

	result, ok := GetIpInfoFromCache(ip)
	assert.False(t, ok)
	assert.Nil(t, result)
}
