package ip_cache_module

import (
	"os"
	"strings"
	"sync"

	app_path_util "github.com/diogopereiradev/httpzen/internal/utils/app_path"
	"github.com/spf13/viper"
)

var (
	CacheFileName      = "ip_cache"
	cacheLock          sync.Mutex
	cacheFileExtension = "json"
	cacheViper         = newCacheViper()
)

func newCacheViper() *viper.Viper {
	v := viper.New()
	v.SetConfigName(CacheFileName)
	v.SetConfigType(cacheFileExtension)
	v.AddConfigPath(app_path_util.GetConfigPath())
	_ = v.ReadInConfig()
	return v
}

func GetIpInfoFromCache(ip string) (map[string]any, bool) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	key := strings.ReplaceAll(ip, ".", "_")
	if cacheViper.IsSet(key) {
		return cacheViper.GetStringMap(key), true
	}
	return nil, false
}

func SetIpInfoToCache(ip string, info map[string]any) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	cacheViper.Set(strings.ReplaceAll(ip, ".", "_"), info)
	_ = cacheViper.WriteConfigAs(app_path_util.GetConfigPath() + "/" + CacheFileName + "." + cacheFileExtension)
}

func ClearCache() {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	_ = cacheViper.WriteConfigAs(app_path_util.GetConfigPath() + "/" + CacheFileName + "." + cacheFileExtension)
	os.Remove(app_path_util.GetConfigPath() + "/" + CacheFileName + "." + cacheFileExtension)
	cacheViper = newCacheViper()
}
