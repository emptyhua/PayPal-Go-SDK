package paypalsdk

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type (
	// CreateBillingPlanResp struct
	CreateBillingPlanResp struct {
		ID                  string              `json:"id,omitempty"`
		State               string              `json:"state,omitempty"`
		PaymentDefinitions  []PaymentDefinition `json:"payment_definitions,omitempty"`
		MerchantPreferences MerchantPreferences `json:"merchant_preferences,omitempty"`
		CreateTime          time.Time           `json:"create_time,omitempty"`
		UpdateTime          time.Time           `json:"update_time,omitempty"`
		Links               []Link              `json:"links,omitempty"`
	}

	// CreateAgreementResp struct
	CreateAgreementResp struct {
		State       string      `json:"state,omitempty"`
		Name        string      `json:"name,omitempty"`
		Description string      `json:"description,omitempty"`
		Plan        BillingPlan `json:"plan,omitempty"`
		Links       []Link      `json:"links,omitempty"`
		StartDate   string      `json:"start_date,omitempty"`
	}
)

func (r CreateAgreementResp) GetExecuteToken() (string, error) {
	for _, link := range r.Links {
		if link.Rel == "approval_url" {
			u, err := url.Parse(link.Href)
			if err != nil {
				return "", err
			}
			q := u.Query()
			return q.Get("token"), nil
		}
	}
	return "", fmt.Errorf("can't find execute token")
}

// CreateBillingPlan creates a billing plan in Paypal
// Endpoint: POST /v1/payments/billing-plans
func (c *Client) CreateBillingPlan(plan BillingPlan) (*CreateBillingPlanResp, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-plans"), plan)
	response := &CreateBillingPlanResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// ActivatePlan activates a billing plan
// By default, a new plan is not activated
// Endpoint: PATCH /v1/payments/billing-plans/
func (c *Client) ActivatePlan(planID string) error {
	buf := bytes.NewBuffer([]byte("[{\"op\":\"replace\",\"path\":\"/\",\"value\":{\"state\":\"ACTIVE\"}}]"))
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-plans/"+planID), buf)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.ClientID, c.Secret)
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)
	return c.SendWithAuth(req, nil)
}

// CreateBillingAgreement creates an agreement for specified plan
// Endpoint: POST /v1/payments/billing-agreements
func (c *Client) CreateBillingAgreement(a BillingAgreement) (*CreateAgreementResp, error) {
	// PayPal needs only ID, so we will remove all fields except Plan ID
	a.Plan = &BillingPlan{
		ID: a.Plan.ID,
	}

	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-agreements"), a)
	response := &CreateAgreementResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)

	return response, err
}

// GetBillingAgreement show agreement details
// Endpoint: GET /v1/payments/billing-agreements/{agreement_id}
func (c *Client) GetBillingAgreement(aid string) (*BillingAgreement, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-agreements/"+aid), nil)
	if err != nil {
		return nil, err
	}
	response := &BillingAgreement{}
	err = c.SendWithAuth(req, response)
	return response, err
}

// ExecuteApprovedAgreement - Use this call to execute (complete) a PayPal agreement that has been approved by the payer.
// Endpoint: POST /v1/payments/billing-agreements/token/agreement-execute
func (c *Client) ExecuteApprovedAgreement(token string) (*BillingAgreement, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-agreements/"+token+"/agreement-execute"), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.ClientID, c.Secret)
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)

	e := BillingAgreement{}

	if err = c.SendWithAuth(req, &e); err != nil {
		return &e, err
	}

	if e.ID == "" {
		return &e, errors.New("Unable to execute agreement with token=" + token)
	}

	return &e, err
}
