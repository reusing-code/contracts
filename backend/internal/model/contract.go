package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Contract struct {
	ID                      uuid.UUID `json:"id"`
	CategoryID              uuid.UUID `json:"categoryId"`
	Name                    string    `json:"name"`
	ProductName             string    `json:"productName,omitempty"`
	Company                 string    `json:"company,omitempty"`
	ContractNumber          string    `json:"contractNumber,omitempty"`
	CustomerNumber          string    `json:"customerNumber,omitempty"`
	PricePerMonth           *float64  `json:"pricePerMonth,omitempty"`
	StartDate               string    `json:"startDate"`
	EndDate                 string    `json:"endDate,omitempty"`
	MinimumDurationMonths   int       `json:"minimumDurationMonths"`
	ExtensionDurationMonths int       `json:"extensionDurationMonths"`
	NoticePeriodMonths      int       `json:"noticePeriodMonths"`
	CustomerPortalURL       string    `json:"customerPortalUrl,omitempty"`
	PaperlessURL            string    `json:"paperlessUrl,omitempty"`
	Comments                string    `json:"comments,omitempty"`
	CreatedAt               time.Time `json:"createdAt"`
	UpdatedAt               time.Time `json:"updatedAt"`
}

type ContractInput struct {
	Name                    string   `json:"name"`
	ProductName             string   `json:"productName,omitempty"`
	Company                 string   `json:"company,omitempty"`
	ContractNumber          string   `json:"contractNumber,omitempty"`
	CustomerNumber          string   `json:"customerNumber,omitempty"`
	PricePerMonth           *float64 `json:"pricePerMonth,omitempty"`
	StartDate               string   `json:"startDate"`
	EndDate                 string   `json:"endDate,omitempty"`
	MinimumDurationMonths   int      `json:"minimumDurationMonths"`
	ExtensionDurationMonths int      `json:"extensionDurationMonths"`
	NoticePeriodMonths      int      `json:"noticePeriodMonths"`
	CustomerPortalURL       string   `json:"customerPortalUrl,omitempty"`
	PaperlessURL            string   `json:"paperlessUrl,omitempty"`
	Comments                string   `json:"comments,omitempty"`
}

func (c *ContractInput) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	if c.StartDate == "" {
		return errors.New("startDate is required")
	}
	return nil
}
