package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/tobi/contracts/backend/internal/middleware"
)

type categorySummary struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	ContractCount int       `json:"contractCount"`
	MonthlyTotal  float64   `json:"monthlyTotal"`
}

type summaryResponse struct {
	TotalContracts    int               `json:"totalContracts"`
	TotalMonthlyAmount float64          `json:"totalMonthlyAmount"`
	Categories        []categorySummary `json:"categories"`
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	cats, err := h.store.ListCategories(r.Context(), middleware.GetUserID(r.Context()))
	if err != nil {
		h.handleStoreError(w, err)
		return
	}

	contracts, err := h.store.ListContracts(r.Context(), middleware.GetUserID(r.Context()))
	if err != nil {
		h.handleStoreError(w, err)
		return
	}

	// Build per-category aggregation
	type agg struct {
		count        int
		monthlyTotal float64
	}
	byCategory := make(map[uuid.UUID]*agg)
	for _, cat := range cats {
		byCategory[cat.ID] = &agg{}
	}

	var totalMonthly float64
	for _, con := range contracts {
		a, ok := byCategory[con.CategoryID]
		if !ok {
			a = &agg{}
			byCategory[con.CategoryID] = a
		}
		a.count++
		if con.PricePerMonth != nil {
			a.monthlyTotal += *con.PricePerMonth
			totalMonthly += *con.PricePerMonth
		}
	}

	catSummaries := make([]categorySummary, 0, len(cats))
	for _, cat := range cats {
		a := byCategory[cat.ID]
		catSummaries = append(catSummaries, categorySummary{
			ID:            cat.ID,
			Name:          cat.Name,
			ContractCount: a.count,
			MonthlyTotal:  a.monthlyTotal,
		})
	}

	h.writeJSON(w, http.StatusOK, summaryResponse{
		TotalContracts:    len(contracts),
		TotalMonthlyAmount: totalMonthly,
		Categories:        catSummaries,
	})
}
