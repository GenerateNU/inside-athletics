import json
import sys
from jwcrypto import jwt, jwk
from datetime import datetime, timedelta


PRIVATE_JWK_DICT = {
    "kty": "EC",
    "kid": "03d4a51d-d8b1-469b-858b-0e586fef312d",
    "use": "sig",
    "key_ops": ["sign", "verify"],
    "alg": "ES256",
    "ext": True,
    "d": "ld8NdnybqZ7xL4w3vOZZ6OAtFt6lqcaCZ_Jno7GF9Co",
    "crv": "P-256",
    "x": "8W4dxm9IaI6k0d4DtgwXtHyxaXfUFsFNROOxLWP_SDI",
    "y": "n1OXBOVJwN64HhGm_sapKD-fhdWUHHggjYnPqbbkofA"
  }

def main():
    key = jwk.JWK(**PRIVATE_JWK_DICT)
    args = sys.argv[1:] 
    header = {
        "alg": "ES256",
        "typ": "JWT",
        "kid": key['kid']
    }

    payload = {
        "sub": args[0],
        "aud": "authenticated", 
        "role": "authenticated",
        "email": "user@example.com",
        "iat": int(datetime.now().timestamp()),
        "exp": int((datetime.now() + timedelta(minutes=15)).timestamp()), 
        "app_metadata": {                           
            "custom_claim": "some_value"
        }
    }

    token = jwt.JWT(header=header, claims=json.dumps(payload))
    token.make_signed_token(key)
    generated_token = token.serialize()
    print("Authorization Bearer:", generated_token)


if __name__ == "__main__":
    main()
