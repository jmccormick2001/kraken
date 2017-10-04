package apiservermsgs

import (
	crv1 "github.com/crunchydata/kraken/apis/cr/v1"
)

type CreatePolicyRequest struct {
	Name       string
	PolicyURL  string
	PolicyFile string
	Namespace  string
}
type ApplyResults struct {
	Results []string
}
type ShowPolicyResponse struct {
	PolicyList crv1.PgpolicyList
}
