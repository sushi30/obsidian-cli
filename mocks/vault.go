package mocks

import (
	"time"

	"github.com/Yakitrak/obsidian-cli/pkg/obsidian"
)

type MockVaultOperator struct {
	DefaultNameErr       error
	PathError            error
	DailyNotePatternErr  error
	Name                 string
	DailyPattern         string
}

func (m *MockVaultOperator) DefaultName() (string, error) {
	if m.DefaultNameErr != nil {
		return "", m.DefaultNameErr
	}
	return m.Name, nil
}

func (m *MockVaultOperator) SetDefaultName(_ string) error {
	if m.DefaultNameErr != nil {
		return m.DefaultNameErr
	}
	return nil
}

func (m *MockVaultOperator) Path() (string, error) {
	if m.PathError != nil {
		return "", m.PathError
	}
	return "path", nil
}

func (m *MockVaultOperator) DailyNotePattern() (string, error) {
	if m.DailyNotePatternErr != nil {
		return "", m.DailyNotePatternErr
	}
	if m.DailyPattern == "" {
		return "", nil
	}
	return m.DailyPattern, nil
}

func (m *MockVaultOperator) ResolveDailyNote() (string, error) {
	pattern, err := m.DailyNotePattern()
	if err != nil {
		return "", err
	}
	if pattern == "" {
		return "", nil
	}
	return obsidian.ExpandDatePattern(pattern, time.Now()), nil
}
