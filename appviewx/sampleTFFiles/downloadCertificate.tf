provider "appviewx"{
  	appviewx_username=""
	appviewx_password=""
	appviewx_environment_is_https=true
	appviewx_environment_ip="192.168.150.192"
	appviewx_environment_port="31443"
}

resource "appviewx_download_certificate" "downloadCert"{
	resource_id="64100a3d5e634315ee9bb536"
     	certificate_download_path="/home/dhivya.v/external_projects/projects/terraform_provider/certs"
     	certificate_download_format="P12"
     	certificate_download_password=""
     	certificate_chain_required=true
     	
     	#Password field is only mandatory for the formats (PFX, P12, JKS)
     	#Private key download access in appviewx is required for PFX, P12, JKS formats
     	#isChainRequired field is only applicable for formats (CRT, CER, CERT, PEM, DER).
}
