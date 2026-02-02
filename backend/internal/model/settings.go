package model

import "time"

type UserSettings struct {
	RenewalDays       int       `json:"renewalDays"`
	ReminderFrequency string    `json:"reminderFrequency"`
	LastReminderSent  time.Time `json:"lastReminderSent,omitempty"`
}

func DefaultUserSettings() UserSettings {
	return UserSettings{
		RenewalDays:       90,
		ReminderFrequency: "disabled",
	}
}

type SettingsResponse struct {
	RenewalDays       int    `json:"renewalDays"`
	ReminderFrequency string `json:"reminderFrequency"`
}
