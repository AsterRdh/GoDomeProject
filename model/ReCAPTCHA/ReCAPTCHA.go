package ReCAPTCHA

type ReCAPTCHATokenResp struct {
	Success     bool     `json:"success"`
	Hostname    string   `json:"challenge_ts"`
	ChallengeTS string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

type ReCAPTCHATokenReq struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
}
