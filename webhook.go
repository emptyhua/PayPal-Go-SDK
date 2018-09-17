package paypalsdk

import (
	"encoding/json"
	"fmt"
	"time"
)

type (
	WebhookEvent struct {
		ID              string          `json:"id"`
		CreateTime      time.Time       `json:"create_time"`
		ResourceType    string          `json:"resource_type"`
		EventVersion    string          `json:"event_version"`
		EventType       string          `json:"event_type"`
		Summary         string          `json:"summary"`
		ResourceVersion string          `json:"resource_version"`
		Resource        json.RawMessage `json:"resource"`
		Links           []Link          `json:"links"`
	}

	WebhookVerifySignatureReq struct {
		AuthAlgo         string       `json:"auth_algo"`
		CertUrl          string       `json:"cert_url"`
		TransmissionId   string       `json:"transmission_id"`
		TransmissionSig  string       `json:"transmission_sig"`
		TransmissionTime string       `json:"transmission_time"`
		WebhookId        string       `json:"webhook_id"`
		WebhookEvent     WebhookEvent `json:"webhook_event"`
	}

	WebhookVerifySignatureResp struct {
		VerificationStatus string `json:"verification_status"`
	}
)

// Verify webhook signature
// Endpoint: POST /v1/notifications/verify-webhook-signature
func (c *Client) VerifyWebhookSignature(params WebhookVerifySignatureReq) (bool, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/verify-webhook-signature"), params)
	response := &WebhookVerifySignatureResp{}
	if err != nil {
		return false, err
	}
	err = c.SendWithAuth(req, response)
	return response.VerificationStatus == "SUCCESS", err
}
