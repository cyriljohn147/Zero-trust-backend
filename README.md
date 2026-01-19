# Zero-Trust Based Secure Backend Authentication System

## 📌 Project Overview
This project implements a **Zero-Trust Based Secure Backend Authentication System** that eliminates implicit trust in devices and requests.  
Every request is continuously verified using **device-bound cryptographic authentication**, **challenge–response mechanisms**, and **short-lived JWTs**.

The system is designed for **modern distributed, cloud, and microservice-based environments**, aligning with Zero-Trust principles recommended by **NIST, IEEE, and industry best practices**.

---

## 🎯 Problem Statement
Traditional backend authentication systems rely heavily on static credentials or long-lived tokens, which are vulnerable to:
- Credential leakage
- Token replay attacks
- Device impersonation
- Unauthorized lateral movement

This project addresses these issues by implementing **device-based cryptographic authentication** where trust is never assumed and is always verified.

---

## 🛡️ Proposed Solution
A **Zero-Trust backend authentication architecture** where:
- Devices authenticate using **public–private key cryptography (ed25519)**
- Authentication uses **challenge–response** instead of passwords
- Access tokens are **short-lived and device-bound**
- Every API request is validated through a **Zero-Trust middleware**
- Compromised devices can be **instantly revoked**
- All security events are **audit logged**

---

## 🧠 Core Features
- One-time device onboarding with public key registration
- Secure challenge–response authentication
- ed25519-based signature verification
- Short-lived JWT issuance
- Zero-Trust request validation middleware
- Device revocation mechanism
- Audit logging for security analysis
- Modular and scalable Go backend architecture

---

## 🏗️ Project Structure

```text
zero-trust-backend/
│
├── cmd/server/           # Application entry point
├── internal/
│   ├── api/              # REST API handlers
│   ├── auth/             # JWT & Zero-Trust middleware
│   ├── crypto/           # ed25519 cryptographic logic
│   ├── db/               # Database access layer
│   ├── services/         # Business logic
│   ├── models/           # Data models
│   └── config/           # Configuration loader
│
├── device-client/        # Sample device-side client
├── migrations/           # Database schema
├── tests/                # Unit & integration tests
├── docs/                 # Architecture & API docs
├── go.mod
└── README.md
```
---

## 🔐 Authentication Flow
1. **Device Registration**
   - Device generates an ed25519 key pair
   - Public key is registered with the backend

2. **Challenge Generation**
   - Backend generates a time-bound random challenge

3. **Challenge Response**
   - Device signs the challenge using its private key

4. **Verification**
   - Backend verifies the signature using the stored public key

5. **Token Issuance**
   - Short-lived JWT is issued and bound to the device

6. **Zero-Trust Enforcement**
   - Every request is validated via middleware
   - Token, device status, and permissions are checked

---

## 🧪 Technology Stack
- **Language:** Go
- **Framework:** Gin
- **Authentication:** ed25519, JWT
- **Database:** PostgreSQL (planned)
- **Architecture:** Modular, layered, Zero-Trust

---

## 📦 Dependencies
| Dependency | Purpose |
|----------|--------|
| gin | REST API & middleware |
| jwt | Short-lived access tokens |
| uuid | Unique device & challenge IDs |
| x/crypto | ed25519 cryptography |

---

## 🚀 Getting Started

### Prerequisites
- Go 1.25+
- PostgreSQL
- OpenSSL (for testing)
- Git

### Environment Setup
1. **Create `.env` file:**
```bash
DATABASE_URL="postgres://user:password@localhost:5432/zerotrust?sslmode=disable"
```

2. **Initialize database:**
```bash
psql -d zerotrust -f migrations/001_init.sql
```

3. **Create a test user:**
```bash
psql -d zerotrust -c "INSERT INTO users (username, email, password_hash, role) VALUES ('testuser', 'test@example.com', 'hash', 'customer');"
```

### Running the Server
```bash
go mod tidy
go run cmd/server/main.go
```

Server starts on `http://localhost:8080`

### API Testing

**1. Register Device**
```bash
curl -X POST http://localhost:8080/devices/register \
  -H "Content-Type: application/json" \
  -d '{"public_key":"YOUR_BASE64_PUBLIC_KEY"}'
```
Response: `{"device_id":"<uuid>"}`

**2. Request Challenge**
```bash
curl -X POST http://localhost:8080/auth/challenge \
  -H "Content-Type: application/json" \
  -d '{"device_id":"<device_id>"}'
```
Response: `{"challenge":"<base64>","challenge_id":"<uuid>","expires_at":"<time>"}`

**3. Verify Challenge (Sign & Authenticate)**
```bash
# Decode challenge
echo "<challenge>" | base64 --decode > challenge.bin

# Sign with private key
openssl pkeyutl -sign -inkey device_private.pem -rawin -in challenge.bin | base64 > signature.txt

# Verify
curl -X POST http://localhost:8080/auth/verify \
  -H "Content-Type: application/json" \
  -d '{"challenge_id":"<challenge_id>","signature":"<signature>"}'
```
Response: `{"access_token":"<jwt>"}`

**4. Test Health Endpoint**
```bash
curl http://localhost:8080/health
```

⸻

## 🔧 Device Client Integration

To automate the authentication process in your application:

```go
type DeviceClient struct {
    privateKey *ed25519.PrivateKey
    deviceID   uuid.UUID
    baseURL    string
}

// Authenticate automatically handles challenge-response
func (dc *DeviceClient) Authenticate() (string, error) {
    // 1. Request challenge
    challenge, challengeID := dc.requestChallenge()
    
    // 2. Sign challenge
    challengeBytes, _ := base64.StdEncoding.DecodeString(challenge)
    signature := ed25519.Sign(dc.privateKey, challengeBytes)
    
    // 3. Verify and get token
    return dc.verifyChallenge(challengeID, base64.StdEncoding.EncodeToString(signature))
}
```

Every API request would then automatically:
1. Call `Authenticate()` to get a fresh token
2. Add `Authorization: Bearer <token>` header
3. Execute the request

⸻

📊 Evaluation Metrics
	•	Authentication latency
	•	Token expiration enforcement
	•	Resistance to replay attacks
	•	Device revocation effectiveness
	•	Challenge-response cryptographic strength

⸻

📚 Academic Relevance
	•	Aligns with Zero-Trust Architecture (NIST SP 800-207)
	•	Uses modern cryptographic practices
	•	Demonstrates real-world security design
	•	Suitable for final-year major project & research publication

⸻

📌 Future Enhancements
	•	Mutual TLS (mTLS)
	•	Hardware-backed key storage
	•	Role-based access control (RBAC)
	•	Rate limiting & anomaly detection
	•	Integration with cloud IAM systems

⸻

🧾 License

This project is developed for academic and research purposes.

⸻
