package apiservermsgs

import (
	crv1 "github.com/crunchydata/kraken/apis/cr/v1"
)

type ShowPolicyResponse struct {
	PolicyList crv1.PgpolicyList
}
