package utils

type UserConfig struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type Config struct {
	ServiceEmail string       `yaml:"service_email"`
	Color        string       `yaml:"color"`
	CompanyName  string       `yaml:"company_name"`
	Users        []UserConfig `yaml:"users"`
}

type ErrorUserNotFound struct{}

func (e *ErrorUserNotFound) Error() string {
	return "User not found"
}

func (conf *Config) GetPasswordHashFromUserEmail(user_email string) (string, error) {
	for _, u := range conf.Users {
		if u.Email == user_email {
			return u.Password, nil
		}
	}

	return "", &ErrorUserNotFound{}
}

// Usage
// import "github.com/ilyakaznacheev/cleanenv"
// type Config struct {
// 	Host         string `yaml:"host" env:"HOST" env-default:"localhost"`
// 	Name         string `yaml:"name" env:"NAME" env-default:"postgres"`
// 	User         string `yaml:"user" env:"USER" env-default:"user"`
// 	Password     string `yaml:"password" env:"PASSWORD"`
// }
// var conf Config
// if err := cleanenv.ReadEnv(&conf); err != nil {
//     ...
// }
