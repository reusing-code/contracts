package model

type UserSettings struct {
	RenewalDays int `json:"renewalDays"`
}

func DefaultUserSettings() UserSettings {
	return UserSettings{RenewalDays: 90}
}
