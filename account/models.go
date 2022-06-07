// Account represents an account in the form3 org section.
// See https://api-docs.form3.tech/api.html#organisation-accounts for
// more information about fields.
package account

type AccountApiResponse struct {
	ResponseBody *Account `json:"response_body,omitempty"`
	StatusCode   *int     `json:"status_code,omitempty"`
	Status       *string  `json:"status,omitempty"`
}

type Account struct {
	Data *AccountData `json:"data,omitempty"`
}

type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`
}

type AccountAttributes struct {
	AccountClassification   *string                `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool                  `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string                 `json:"account_number,omitempty"`
	AlternativeNames        []string               `json:"alternative_names,omitempty"`
	BankID                  string                 `json:"bank_id,omitempty"`
	BankIDCode              string                 `json:"bank_id_code,omitempty"`
	BaseCurrency            string                 `json:"base_currency,omitempty"`
	Bic                     string                 `json:"bic,omitempty"`
	Country                 *string                `json:"country,omitempty"`
	Iban                    string                 `json:"iban,omitempty"`
	JointAccount            *bool                  `json:"joint_account,omitempty"`
	Name                    []string               `json:"name,omitempty"`
	SecondaryIdentification string                 `json:"secondary_identification,omitempty"`
	Status                  *string                `json:"status,omitempty"`
	Switched                *bool                  `json:"switched,omitempty"`
	CustomerID              *string                `json:"customer_id,omitempty"`
	NameMatchingStatus      *string                `json:"name_matching_status,omitempty"`
	StatusReason            *string                `json:"status_reason,omitempty"`
	PrivateIdentification   *PrivateIdentification `json:"private_identification,omitempty"`
	UserDefinedData         *[]UserDefinedData     `json:"user_defined_data,omitempty"`
	ValidationType          *string                `json:"validation_type,omitempty"`
	ReferenceMask           string                 `json:"reference_mask,omitempty"`
	AcceptanceQualifier     string                 `json:"acceptance_qualifier,omitempty"`
	Relationships           *Relationships         `json:"relationships,omitempty"`
}

type Relationships struct {
	AccountEvents *AccountEvents `json:"account_events,omitempty"`
	MasterAccount *MasterAccount `json:"master_account,omitempty"`
}

type MasterAccount struct {
	Data *[]MasterAccountData `json:"data,omitempty"`
}

type MasterAccountData struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type AccountEvents struct {
	Data *[]AccountEventsData `json:"data,omitempty"`
}

type AccountEventsData struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type UserDefinedData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PrivateIdentification struct {
	Identification string   `json:"identification,omitempty"`
	BirthDate      string   `json:"birth_date,omitempty"`
	BirthCountry   string   `json:"birth_country,omitempty"`
	Address        []string `json:"address,omitempty"`
	City           string   `json:"city,omitempty"`
	Country        string   `json:"country,omitempty"`
}

type OrganizationIdentification struct {
	Identification string   `json:"identification,omitempty"`
	Actors         *[]Actor `json:"actors,omitempty"`
	Address        string   `json:"address,omitempty"`
	City           string   `json:"city,omitempty"`
	Country        string   `json:"country,omitempty"`
}

type Actor struct {
	Name      string `json:"name,omitempty"`
	BirthDate string `json:"birth_date,omitempty"`
	Residency string `json:"residency,omitempty"`
}
