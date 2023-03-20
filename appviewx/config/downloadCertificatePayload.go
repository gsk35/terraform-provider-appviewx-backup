package config

type DownloadCertificatePayload struct {
	CommonName      string `json:"commonName"`
	SerialNumber    string `json:"serialNumber"`
	Format          string `json:"format"`
	Password        string `json:"password"`
	IsChainRequired bool   `json:"isChainRequired"`
}

type SearchCertificatePayload struct {
	Input  Input  `json:"input"`
	Filter Filter `json:"filter"`
}

type Input struct {
	ResourceId string `json:"resourceId"`
}

type Filter struct {
	SortOrder string `json:"sortOrder"`
}

type AppviewxSearchCertResponse struct {
	AppviewxResponse CertificateResponse `json:"response"`
	Message          string              `json:"message"`
	AppStatusCode    string              `json:"appStatusCode"`
	Tags             map[string]string   `json:"tags"`
	Headers          string              `json:"headers"`
}

type CertificateResponse struct {
	ResponseObject CertificateObjects `json:"response"`
}

type CertificateObjects struct {
	Objects []CertificateDetails `json:"objects"`
}

type CertificateDetails struct {
	CommonName   string `json:"commonName"`
	SerialNumber string `json:"serialNumber"`
}

var DownloadCertificateActionId = "certificate/download/format"
var SearchCertificateActionId = "certificate/search"
