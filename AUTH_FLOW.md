# Zero Trust Authentication Flow

Complete step-by-step guide to get an access token and test the secure endpoint.

## Prerequisites

- Ed25519 keypair (`private_key.pem` and `public_key.pem`)
- Server running on `http://localhost:8080`
- PostgreSQL database initialized

## Step 1: Generate Ed25519 Keypair

```bash
# Create a new private key
openssl genpkey -algorithm ED25519 -out private_key.pem

# Extract the public key
openssl pkey -in private_key.pem -pubout -out public_key.pem
```

## Step 2: Register Device

Extract the raw 32-byte Ed25519 public key:

```bash
PUBLIC_KEY=$(openssl pkey -in private_key.pem -pubout -outform DER | tail -c 32 | base64 | tr -d '\n')

curl -X POST http://localhost:8080/devices/register \
  -H "Content-Type: application/json" \
  -d "{\"user_id\": 1, \"public_key\": \"$PUBLIC_KEY\"}"
```

**Response:**
```json
{"device_id":"5fc379f5-8fbb-42d4-a17e-081eabcacdd6"}
```

Save the `device_id`.

## Step 3: Request Challenge

```bash
DEVICE_ID="5fc379f5-8fbb-42d4-a17e-081eabcacdd6"

curl -X POST http://localhost:8080/auth/challenge \
  -H "Content-Type: application/json" \
  -d "{\"device_id\": \"$DEVICE_ID\"}"
```

**Response:**
```json
{
  "challenge":"EuRj1+vOQTBoKtxzCvOlo+wP+QZxI+V+2SWgKteKpb4=",
  "challenge_id":"77ca81a9-875f-4830-9757-ac0a98109a8f",
  "expires_at":"2026-01-19T17:48:28.595554+05:30"
}
```

Save `challenge` and `challenge_id`.

## Step 4: Sign the Challenge

Option A: Use Python script Manual:

```bash
python3 sign_challenge.py "$CHALLENGE"
```

Option B: Manual process(Recommended):

```bash
CHALLENGE="EuRj1+vOQTBoKtxzCvOlo+wP+QZxI+V+2SWgKteKpb4="

curl -X POST http://localhost:8080/auth//sign \
  -H "Content-Type: application/json" \
  -d "{\"challenge\": \"$CHALLENGE\"}"
```

**Response:**
```json
{
  "signature":"rwR+c9dXBFR2yewepnZjES7TepvFWJHCeULf1nB/Cx4I5qRQzn+85y+wA0vcqkK+jOvf8SwSWdTGVrdKHNBYCA=="
}
```

**Output:** Base64-encoded signature

## Step 5: Verify Challenge & Get Token

```bash
SIGNATURE="rwR+c9dXBFR2yewepnZjES7TepvFWJHCeULf1nB/Cx4I5qRQzn+85y+wA0vcqkK+jOvf8SwSWdTGVrdKHNBYCA=="

curl -X POST http://localhost:8080/auth/verify \
  -H "Content-Type: application/json" \
  -d "{
    \"challenge_id\": \"$CHALLENGE_ID\",
    \"signature\": \"$SIGNATURE\"
  }"
```

**Response:**
```json
{
  "access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkZXZpY2VfaWQiOiI1ZmMzNzlmNS04ZmJiLTQyZDQtYTE3ZS0wODFlYWJjYWNkZDYiLCJ1c2VyX2lkIjoxLCJleHAiOjE3Njg4MjUzMjIsImlhdCI6MTc2ODgyNTAyMn0.bhrK0VpHMh00ISl6jmBaijP7p_F2dP8CPN9fpICIAU0"
}
```

## Step 6: Access Protected Endpoint

```bash
ACCESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkZXZpY2VfaWQiOiI1ZmMzNzlmNS04ZmJiLTQyZDQtYTE3ZS0wODFlYWJjYWNkZDYiLCJ1c2VyX2lkIjoxLCJleHAiOjE3Njg4MjUzMjIsImlhdCI6MTc2ODgyNTAyMn0.bhrK0VpHMh00ISl6jmBaijP7p_F2dP8CPN9fpICIAU0"

curl http://localhost:8080/api/secure-data \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Response:**
```json
{"message":"Zero Trust access granted"}
```

## Test Device Revocation

```bash
DEVICE_ID="5fc379f5-8fbb-42d4-a17e-081eabcacdd6"

# Revoke the device
psql -U cyriljohn147 -d zerotrust -c "UPDATE devices SET status='revoked' WHERE device_id='$DEVICE_ID';"

# Try accessing with the same token
curl http://localhost:8080/api/secure-data \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Result:** `403 Forbidden` - Device instantly blocked!

## Key Points

- ✅ Token expires in 5 minutes
- ✅ Device status checked on every request
- ✅ Challenge is valid for 2 minutes
- ✅ All keys must be base64-encoded
- ✅ Public key registered: raw 32-byte Ed25519 format (extracted from DER)
- ✅ Instant revocation with no cache delays
