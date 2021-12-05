//go:build !sqlite
// +build !sqlite

package driver

func (m *RegistrySQL) CanHandle(dsn string) bool {
	return m.alwaysCanHandle(dsn)
}
