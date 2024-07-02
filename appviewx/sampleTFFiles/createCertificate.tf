provider "appviewx"{
  	appviewx_username=""
	appviewx_password=""
	appviewx_environment_is_https=true
	appviewx_environment_ip=""
	appviewx_environment_port=""
}

resource "appviewx_create_certificate" "sampletest"{
	common_name="msca.appviewx.com"
	hash_function="SHA256"
      	key_type="RSA"
     	bit_length="2048"
     	certificate_authority="Microsoft Enterprise"
     	ca_setting_name="microsoft enterprise"
     	#dns_names=["test.com","test123.com"]
     	validity_days=365
     	custom_fields={"test":"msca"}
     	vendor_specific_fields={"templateName":"template"}
     	is_sync=true
     	
     	#if is_sync is true
     	certificate_download_path="/home/certs"
     	certificate_download_format="PEM"
     	certificate_download_password=""
     	certificate_chain_required=true
     	
     	#Password field is only mandatory for the formats (PFX, P12, JKS)
     	#Private key download access in appviewx is required for PFX, P12, JKS formats
     	#isChainRequired field is only applicable for formats (CRT, CER, CERT, PEM, DER).
}
