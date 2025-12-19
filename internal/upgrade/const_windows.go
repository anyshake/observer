package upgrade

const (
	VERSION_CHECK_DOMAIN       = "observer.updates.anyshake.org"
	RELEASE_EXECUTABLE_NAME    = "observer.exe"
	RELEASE_FETCH_URL_TEMPLATE = "https://github.com/anyshake/observer/releases/download/{{.Version}}/{{.ToolchainName}}.{{.Extension}}"
)
