package config

const DEV = "dev"

// Имена переменных окружения
const (
	DbUser                  = "DB_USER_NAME"
	DbPassword              = "DB_USER_PASSWORD"
	DbPort                  = "DB_PORT"
	DbHost                  = "DB_HOST"
	DbName                  = "DB_NAME"
	RootPath                = "ROOT_PATH"
	AdminPassword           = "ADMIN_PASSWORD"
	SessionSecret           = "SESSION_SECRET"
	DeleteTablesBeforeStart = "DELETE_TABLES_BEFORE_START"
	MODE                    = "MODE"
	RootPathWeb             = "ROOT_PATH_WEB"
	SessionMaxAge           = "SESSION_MAX_AGE"
	Salt                    = "SALT"
	PasswordCost            = "PASSWORD_COST"
)

// Configuration Структура для хранения конфигруации
type Configuration struct {
	DbUser                  string `json:"DB_USER,omitempty"`
	DbPassword              string `json:"DB_PASSWORD,omitempty"`
	DbHost                  string `json:"DB_HOST,omitempty"`
	DbPort                  string `json:"DB_PORT,omitempty"`
	DbName                  string `json:"DB_NAME,omitempty"`
	RootPath                string `json:"ROOT_PATH,omitempty"`
	AdminPassword           string `json:"ADMIN_PASSWORD,omitempty"`
	SessionsSecret          string `json:"SESSIONS_SECRET,omitempty"`
	DeleteTablesBeforeStart int    `json:"DELETE_TABLES_BEFORE_START,omitempty"`
	MODE                    string `json:"MODE,omitempty"`
	RootPathWeb             string `json:"ROOT_PATH_WEB,omitempty"`
	SessionMaxAge           int    `json:"SESSION_MAX_AGE,omitempty"`
	Salt                    string `json:"salt"`
	PasswordCost            int    `json:"password_cost"`
}
