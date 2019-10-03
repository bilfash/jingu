package config_test

import (
	"github.com/bilfash/jingu/config"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestConfigGetHostNameShouldReturnHostname(t *testing.T) {
	config := config.New("example.yml")
	actual := config.Host()
	assert.Equal(t, actual, "imap.google.com")
}

func TestConfigGetPortShouldReturnPort(t *testing.T) {
	config := config.New("example.yml")
	actual := config.Port()
	assert.Equal(t, actual, "993")
}

func TestConfigGetUsernameShouldReturnUsername(t *testing.T) {
	config := config.New("example.yml")
	actual := config.Username()
	assert.Equal(t, actual, "jingu@gmail.com")
}

func TestConfigGetPasswordShouldReturnPassword(t *testing.T) {
	config := config.New("example.yml")
	actual := config.Password()
	assert.Equal(t, actual, "jingu")
}

func TestConfigGetSubjectsShouldReturnSubjects(t *testing.T) {
	config := config.New("example.yml")
	actual := config.Subjects()
	assert.Equal(t, actual, []string{"subject1", "subject2", "subject3"})
}

func TestConfigGetFilePatternShouldReturnFilePattern(t *testing.T) {
	config := config.New("example.yml")
	actual := config.FilePattern()
	assert.Equal(t, actual, ".txt$")
}

func TestConfigGetSinkFolderShouldReturnSinkFolder(t *testing.T) {
	config := config.New("example.yml")
	actual := config.SinkFolder()
	assert.Equal(t, actual, "./")
}
