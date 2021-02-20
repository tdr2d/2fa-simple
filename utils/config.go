package utils

type UserConfig struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type Config struct {
	BaseUrl      string       `yaml:"base_url"`
	TemplateDir  string       `yaml:"template_dir"`
	SpaDirectory string       `yaml:"spa_directory"`
	SpaFallback  string       `yaml:"spa_fallback"`
	ServiceEmail string       `yaml:"service_email"`
	Color        string       `yaml:"color"`
	CompanyName  string       `yaml:"company_name"`
	Users        []UserConfig `yaml:"users"`
}

type ErrorUserNotFound struct{}

func (e *ErrorUserNotFound) Error() string {
	return "user_not_found"
}

func (conf *Config) GetPasswordHashFromUserEmail(user_email string) (string, error) {
	for _, u := range conf.Users {
		if u.Email == user_email {
			return u.Password, nil
		}
	}

	return "", &ErrorUserNotFound{}
}
