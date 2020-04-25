package envconfig

const (
	keyUserUid      string = "AS_USER_UID"
	keyUserGid      string = "AS_USER_GID"
	keyRootRequired string = "AS_ROOT_REQUIRED"
)

func GetDefaultConfigValues() EnvConfig {
	return EnvConfig{
		keyUserUid:      "1000",
		keyUserGid:      "1000",
		keyRootRequired: "false",
	}
}
