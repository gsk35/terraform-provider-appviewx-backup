provider "appviewx"{
  	appviewx_username=""
	appviewx_password=""
	appviewx_environment_is_https=true
	appviewx_environment_ip=""
	appviewx_environment_port=""
}

resource "appviewx_download_certificate" "downloadCert"{
	resource_id=""
     	certificate_download_path="/home/certs"
     	certificate_download_format="P12"
     	certificate_download_password=""
     	certificate_chain_required=true
     	
     	#Password field is only mandatory for the formats (PFX, P12, JKS)
     	#Private key download access in appviewx is required for PFX, P12, JKS formats
     	#isChainRequired field is only applicable for formats (CRT, CER, CERT, PEM, DER).
}
