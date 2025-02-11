//go:build dev

package config

func init() {
	ENV = "dev"
	ZIPCLOUD_API_URL = "https://zipcloud.ibsnet.co.jp/api/search"
}
