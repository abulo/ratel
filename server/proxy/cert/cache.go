package cert

import "crypto/tls"

// Cache 证书缓存接口
type Cache interface {
	Set(host string, c *tls.Certificate)
	Get(host string) *tls.Certificate
}
