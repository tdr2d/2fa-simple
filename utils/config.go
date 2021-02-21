package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type UserConfig struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type Config struct {
	BaseUrl        string `yaml:"base_url"`
	TemplateDir    string `yaml:"template_dir"`
	SpaDirectory   string `yaml:"spa_directory"`
	SpaFallback    string `yaml:"spa_fallback"`
	SqliteDatabase string `yaml:"sqlite_database"`

	Website      string `yaml:"website"`
	ServiceEmail string `yaml:"service_email"`
	Color        string `yaml:"color"`
	CompanyName  string `yaml:"company_name"`

	Users []UserConfig `yaml:"users"`
}

type ErrorUserNotFound struct{}

func (e *ErrorUserNotFound) Error() string {
	return "user_not_found"
}

func (conf *Config) UserExists(user_email string) bool {
	for _, u := range conf.Users {
		if u.Email == user_email {
			return true
		}
	}
	return false
}

func (conf *Config) GetPasswordHashFromUserEmail(user_email string) (string, error) {
	for _, u := range conf.Users {
		if u.Email == user_email {
			return u.Password, nil
		}
	}
	return "", &ErrorUserNotFound{}
}

func (conf *Config) ChangePassword(user_email string, password_hash string) error {
	for i, u := range conf.Users {
		if u.Email == user_email {
			conf.Users[i].Password = password_hash
			return nil
		}
	}
	return &ErrorUserNotFound{}
}

func (conf *Config) EnsureFilesExist() error {
	dir := filepath.Dir(conf.SqliteDatabase)
	if dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	return TouchFile(conf.SqliteDatabase)
}

func TouchFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func WriteYaml(in interface{}, outputfile string) error {
	out, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	if err := TouchFile(outputfile); err != nil {
		return err
	}
	if err := ioutil.WriteFile(outputfile, out, 0644); err != nil {
		return err
	}
	return nil
}
