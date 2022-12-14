package controller

import (
	"awesomeProject1/models"
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/url"
	"os"
)

type response struct {
	CompanyName       string  `json:"company_name"`
	Npwp              string  `json:"npwp"`
	BankDebtToEquity  float64 `json:"bank_debt_to_equity"`
	Capitalisation    float64 `json:"capitalisation"`
	GrossProfitMargin float64 `json:"gross_profit_margin"`
	CurrentRatio      float64 `json:"current_ratio"`
}

type PbkDumps struct {
	CompanyName string `json:"company_name"`
	Npwp        string `json:"npwp"`
	ApiKey      string `json:"api_key"`
}

type Payload struct {
	CompanyName string `json:"company_name"`
	Npwp        string `json:"npwp"`
}

func (request PbkDumps) validate() url.Values {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	errs := url.Values{}

	if request.CompanyName == "" {
		errs.Add("company_name", "Key is required")
	}

	if request.Npwp == "" {
		errs.Add("npwp", "Key is required")
	}

	if request.ApiKey != os.Getenv("API_KEY") {
		errs.Add("api_key", "api_key is wrong")
	}

	return errs
}

func PostCompany(w http.ResponseWriter, r *http.Request) {

	PbkDumps := &PbkDumps{}

	err := json.NewDecoder(r.Body).Decode(&PbkDumps)
	if err != nil {
		log.Fatalf("Can't decode from request body.  %v", err)
	}

	if validErrs := PbkDumps.validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "appliciation/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	Payload := Payload{
		CompanyName: PbkDumps.CompanyName,
		Npwp:        PbkDumps.Npwp,
	}

	data := models.GetCompany(models.Payload(Payload))

	res := response{
		CompanyName:       data.CompanyName,
		Npwp:              data.Npwp,
		BankDebtToEquity:  data.BankDebtToEquity,
		Capitalisation:    data.Capitalisation,
		GrossProfitMargin: data.GrossProfitMargin,
		CurrentRatio:      data.CurrentRatio,
	}

	_ = json.NewEncoder(w).Encode(res)
}
