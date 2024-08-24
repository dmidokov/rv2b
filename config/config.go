// Пакет config осуществляет загрузку конфигурации приложения
// из переменных окружения
package config

import (
	"log"
	"os"
	"strconv"
)

// LoadConfig Функция заполняет структуру config из переменных окружения
// возвращает заполненную структуру и ошибку
//
// Примечание: возможно имеет смысл создать слайс
// с именами переменных которые требуются и далее загрузить
// их в цикле без дублирования кода
func LoadConfig() *Configuration {

	var config = Configuration{}
	var exist bool
	var err error

	if config.DbUser, exist = os.LookupEnv(DbUser); !exist {
		log.Fatalf("database  username is missing")
	}

	if config.DbPassword, exist = os.LookupEnv(DbPassword); !exist {
		log.Fatalf("database password is missing")
	}

	if config.DbHost, exist = os.LookupEnv(DbHost); !exist {
		log.Fatal("database host is missing")
	}

	if config.DbPort, exist = os.LookupEnv(DbPort); !exist {
		log.Fatal("database port is missing")
	}

	if config.DbName, exist = os.LookupEnv(DbName); !exist {
		log.Fatal("database name is missing")
	}

	if config.RootPath, exist = os.LookupEnv(RootPath); !exist {
		log.Fatal("application root path is empty")
	}

	if config.AdminPassword, exist = os.LookupEnv(AdminPassword); !exist {
		log.Fatal("admin password is missing")
	}

	if config.SessionsSecret, exist = os.LookupEnv(SessionSecret); !exist {
		log.Fatal("session secrec is empty")
	}

	if value, exist := os.LookupEnv(DeleteTablesBeforeStart); !exist {
		log.Fatal("delete tables key is empty")
	} else {
		config.DeleteTablesBeforeStart, err = strconv.Atoi(value)
		if err != nil {
			log.Fatal("delete tables key is wrong")
		}
	}

	if config.MODE, exist = os.LookupEnv(MODE); !exist || !(config.MODE == "dev" || config.MODE == "production" || config.MODE == "mock") {
		log.Fatal("MODE is not exist")
	}

	if config.RootPathWeb, exist = os.LookupEnv(RootPathWeb); !exist {
		log.Fatal("application web root path is empty")
	}

	if value, exist := os.LookupEnv(SessionMaxAge); !exist {
		config.SessionMaxAge = 3600
	} else {
		config.SessionMaxAge, err = strconv.Atoi(value)
		if err != nil {
			log.Fatal("session max age key is wrong")
		}
	}

	if config.Salt, exist = os.LookupEnv(Salt); !exist {
		log.Fatal("application web root path is empty")
	}

	if value, exist := os.LookupEnv(PasswordCost); !exist {
		config.PasswordCost = 14
		log.Fatal("application web root path is empty")
	} else {
		config.PasswordCost, err = strconv.Atoi(value)
		if err != nil {
			log.Fatal("password cost key is wrong")
		}
	}

	if config.MigrationPath, exist = os.LookupEnv(MigrationPath); !exist {
		log.Fatal("migration path is empty")
	}

	if config.HttpPort, exist = os.LookupEnv(HttpPort); !exist {
		log.Fatal("http port is empty")
	}

	if config.SSLPort, exist = os.LookupEnv(SSLPort); !exist {
		log.Fatal("ssl port is empty")
	}

	if config.SecretsPath, exist = os.LookupEnv(SecretsPath); !exist {
		log.Fatal("secrets path is empty")
	}

	if config.TempFolder, exist = os.LookupEnv(TempFolder); !exist {
		config.TempFolder = "/bin/temp/"
	}

	return &config
}
