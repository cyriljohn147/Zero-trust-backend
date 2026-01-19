#!/usr/bin/env python3
"""
Sign challenge for Zero Trust authentication.

Usage:
    python3 sign_challenge.py "base64-challenge-string"
    
Output:
    Base64-encoded signature to use in /auth/verify
"""

import base64
import sys
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.backends import default_backend


def sign_challenge(challenge_b64, private_key_path="private_key.pem"):
    """
    Sign a challenge with an Ed25519 private key.
    
    Args:
        challenge_b64: Base64-encoded challenge string
        private_key_path: Path to private_key.pem
        
    Returns:
        Base64-encoded signature
    """
    try:
        # Load private key
        with open(private_key_path, 'rb') as f:
            private_key = serialization.load_pem_private_key(
                f.read(),
                password=None,
                backend=default_backend()
            )
        
        # Decode challenge
        challenge_bytes = base64.b64decode(challenge_b64)
        
        # Sign
        signature = private_key.sign(challenge_bytes)
        
        # Return base64
        return base64.b64encode(signature).decode()
        
    except FileNotFoundError:
        print(f"Error: {private_key_path} not found")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python3 sign_challenge.py '<base64-challenge>'")
        print("\nExample:")
        print("  python3 sign_challenge.py 'EuRj1+vOQTBoKtxzCvOlo+wP+QZxI+V+2SWgKteKpb4='")
        sys.exit(1)
    
    challenge = sys.argv[1]
    signature = sign_challenge(challenge)
    print(signature)
