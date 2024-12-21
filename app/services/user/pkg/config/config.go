package config

import (
	"errors"
	"flag"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ServiceConfig struct {
	SrvConf  ServerConfig
	DBConf   DBConfig
	MailConf ServerMailConf
	Env      string `yaml:"env" env-default:"local"`
}

// кофигурация базы данных
type DBConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
}

// конфигурация сервера
type ServerConfig struct {
	Port    string        `yaml:"port"`
	Host    string        `yaml:"host"`
	Timeout time.Duration `yaml:"timeout"`
}

type ServerMailConf struct {
	Login    string
	Password string
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

// конфигурация REST API Сервера

// MustLoadEnv загружает переменные окружения из файла .env,
// возвращает установленное окружение (local/dev/prod)
func MustLoadEnv() string {
	fi := "config.MustLoadEnv"

	//путь до файла .env
	if err := godotenv.Load("../.env"); err != nil {
		panic(fi + ": " + err.Error())
	}

	env := os.Getenv("ENVIRONMENT")

	return env
}

func MustLoadConfig() ServiceConfig {
	fi := "config.MustLoadConfig"

	//путь до файла конфигурации
	pathToConfDir, nameOfConfFile, err := getConfigLocation()
	if err != nil {
		panic(fi + ": " + err.Error())
	}

	//проверяем существует ли такие директория и файл
	if _, err := os.Stat(pathToConfDir + "/" + nameOfConfFile); os.IsNotExist(err) {
		panic(fi + ": " + err.Error())
	}

	//загружаем конфигурацию
	UserConf, err := LoadConfig(pathToConfDir, nameOfConfFile)
	if err != nil {
		panic(fi + ": " + err.Error())
	}

	return *UserConf

}

func getConfigLocation() (string, string, error) {
	fi := "config.getConfigLocation"

	//загрузка пути к директории с файлами конфигурции и имени файла из argv
	pathToConfDir, nameOfConfFile := getConfLocationFromArgv()

	//если имя директории - пустая строка, пробуем взять его из переменных окружения
	if pathToConfDir == "" {

		pathToConfDir = os.Getenv("CONFIG_DIR")

		if pathToConfDir == "" {
			return "", "", errors.New(fi + ": " + "pathToConfDir is empty at argv and env")
		}
	}

	//если имя файла - пустая строка, пробуем взять его из переменных окружения
	if nameOfConfFile == "" {

		nameOfConfFile = os.Getenv("CONFIG_FILE")

		if nameOfConfFile == "" {
			return "", "", errors.New(fi + ": " + "nameOfConfFile is empty at argv and env")
		}
	}

	return pathToConfDir, nameOfConfFile, nil
}

func getConfLocationFromArgv() (string, string) {

	var (
		pathToConfDir  string
		nameOfConfFile string
	)

	flag.StringVar(&pathToConfDir, "config_path", "", "name of directory with configs")
	flag.StringVar(&nameOfConfFile, "config_file", "", "name of config file")
	flag.Parse()

	return pathToConfDir, nameOfConfFile
}

func LoadConfig(path string, name string) (*ServiceConfig, error) {
	fi := "config.LoadConfig"

	var (
		dbConf   DBConfig
		srvConf  ServerConfig
		mailConf ServerMailConf
	)

	//инициализируем имя, папку и тип конфига
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.New(fi + ": " + err.Error())
	}

	//заполняем структуру сервера отправки писем
	if err := viper.UnmarshalKey("mail", &mailConf); err != nil {
		return nil, err
	}

	//получаем пароль и логин для почты, с которой будем посылать письма
	if mailConf.Login, mailConf.Password = getMailCredentials(); mailConf.Password == "" || mailConf.Login == "" {
		return nil, errors.New("mail credentials are empty")
	}

	//заполняем структуру ДБ
	if err := viper.UnmarshalKey("db", &dbConf); err != nil {
		return nil, err
	}

	//заполняем структуру сервера
	if err := viper.UnmarshalKey("server", &srvConf); err != nil {
		return nil, err
	}

	return &ServiceConfig{
		SrvConf:  srvConf,
		DBConf:   dbConf,
		MailConf: mailConf,
	}, nil

}

func getMailCredentials() (string, string) {

	login := os.Getenv("MAIL_LOGIN")
	password := os.Getenv("MAIL_PASSWORD")

	return login, password
}
