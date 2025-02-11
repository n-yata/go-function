//go:build prd

package config

func init() {
	ENV = "prd"
	ZIPCLOUD_API_URL = "https://zipcloud.ibsnet.co.jp/api/search"
}
