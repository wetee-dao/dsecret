package types

type ReencryptSecretRequest struct {
	// 来源ID
	OrgId string `json:"org_id,omitempty"`
	// 密文ID
	SecretId string `json:"secret_id,omitempty"`
	// 密文接收者公钥
	RdrPk *PubKey `json:"rdr_pk,omitempty"`
}

type ReencryptedSecretShare struct {
	// 来源ID
	OrgId string `json:"org_id,omitempty"`
	// 密文ID
	SecretId string `json:"secret_id,omitempty"`
	// 密文接收者公钥
	RdrPk *PubKey `json:"rdr_pk,omitempty"`
	// 密钥碎片索引
	Index int32 `json:"index,omitempty"`
	// Re-encrypted secret share
	// 重新加密的秘密份额
	XncSki []byte `json:"xnc_ski,omitempty"`
	// Random oracle challenge
	// 随机神谕挑战
	Chlgi []byte `json:"chlgi,omitempty"`
	// NIZK proofi of re-encryption
	// NIZK 重新加密证明
	Proofi []byte `json:"proofi,omitempty"`
}
