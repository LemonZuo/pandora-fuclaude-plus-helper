package util

import (
	"time"
)

type Response struct {
	Accounts        map[string]AccountInfo `json:"accounts"`
	AccountOrdering []string               `json:"account_ordering"`
}

type AccountInfo struct {
	Account                 Account       `json:"account"`
	Features                []string      `json:"features"`
	Entitlement             Entitlement   `json:"entitlement"`
	RateLimits              []interface{} `json:"rate_limits"`
	LastActiveSubscription  Subscription  `json:"last_active_subscription"`
	IsEligibleForYearlyPlus bool          `json:"is_eligible_for_yearly_plus_subscription"`
	CanAccessWithSession    bool          `json:"can_access_with_session"`
	SSOConnectionName       interface{}   `json:"sso_connection_name"`
}

type Account struct {
	AccountUserRole                       string      `json:"account_user_role"`
	AccountUserID                         string      `json:"account_user_id"`
	AccountResidencyRegion                string      `json:"account_residency_region"`
	Processor                             Processor   `json:"processor"`
	AccountID                             string      `json:"account_id"`
	OrganizationID                        string      `json:"organization_id"`
	IsMostRecentExpiredSubscriptionGratis bool        `json:"is_most_recent_expired_subscription_gratis"`
	HasPreviouslyPaidSubscription         bool        `json:"has_previously_paid_subscription"`
	Name                                  string      `json:"name"`
	ProfilePictureID                      interface{} `json:"profile_picture_id"`
	ProfilePictureURL                     interface{} `json:"profile_picture_url"`
	Structure                             string      `json:"structure"`
	PlanType                              string      `json:"plan_type"`
	IsDeactivated                         bool        `json:"is_deactivated"`
	PromoData                             struct{}    `json:"promo_data"`
	ResellerHostedAccount                 bool        `json:"reseller_hosted_account"`
	ResellerID                            interface{} `json:"reseller_id"`
}

type Processor struct {
	A001 ProcessorInfo `json:"a001"`
	B001 ProcessorInfo `json:"b001"`
	C001 ProcessorInfo `json:"c001"`
}

type ProcessorInfo struct {
	HasCustomerObject     bool `json:"has_customer_object,omitempty"`
	HasTransactionHistory bool `json:"has_transaction_history,omitempty"`
}

type Entitlement struct {
	SubscriptionID        string     `json:"subscription_id"`
	HasActiveSubscription bool       `json:"has_active_subscription"`
	SubscriptionPlan      string     `json:"subscription_plan"`
	ExpiresAt             *time.Time `json:"expires_at"`
	BillingPeriod         string     `json:"billing_period"`
}

type Subscription struct {
	SubscriptionID         string `json:"subscription_id"`
	PurchaseOriginPlatform string `json:"purchase_origin_platform"`
	WillRenew              bool   `json:"will_renew"`
}
