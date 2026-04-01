"use client";

import { useEffect, useState } from "react";
import {
  registerDevice,
  getChallenge,
  verifyChallenge,
  getSecureData,
  revokeDevice,
} from "../services/api";
import { generateKeyPair, signChallenge } from "../utils/crypto";

export default function AuthFlow() {
  const [publicKey, setPublicKey] = useState("");
  const [privateKey, setPrivateKey] = useState("");
  const [deviceId, setDeviceId] = useState("");
  const [challenge, setChallenge] = useState("");
  const [challengeId, setChallengeId] = useState("");
  const [signature, setSignature] = useState("");
  const [token, setToken] = useState("");
  const [response, setResponse] = useState("");
  const [revokeResponse, setRevokeResponse] = useState("");
  const [authStatus, setAuthStatus] = useState("");
  const [loading, setLoading] = useState(false);

  // 🔐 Generate keys
  useEffect(() => {
    try {
      const storedPrivate = localStorage.getItem("privateKey");
      const storedPublic = localStorage.getItem("publicKey");

      if (!storedPrivate || !storedPublic) {
        const keys = generateKeyPair();
        localStorage.setItem("privateKey", keys.privateKey);
        localStorage.setItem("publicKey", keys.publicKey);

        setPrivateKey(keys.privateKey);
        setPublicKey(keys.publicKey);
      } else {
        setPrivateKey(storedPrivate);
        setPublicKey(storedPublic);
      }
    } catch (error) {
      console.error("Key error:", error);
    }
  }, []);

  const handleRegister = async () => {
    setLoading(true);
    try {
      const res = await registerDevice(publicKey);
      setDeviceId(res.data.device_id);

      // reset state
      setToken("");
      setResponse("");
      setSignature("");
      setChallenge("");
      setAuthStatus("✅ Device registered");
    } catch {
      setAuthStatus("❌ Registration failed");
    }
    setLoading(false);
  };

  const handleChallenge = async () => {
    setLoading(true);
    try {
      const res = await getChallenge(deviceId);
      setChallenge(res.data.challenge);
      setChallengeId(res.data.challenge_id);
      setAuthStatus("📩 Challenge received");
    } catch {
      setAuthStatus("❌ Failed to get challenge");
    }
    setLoading(false);
  };

  const handleSign = () => {
    try {
      const sig = signChallenge(challenge, privateKey);
      setSignature(sig);
      setAuthStatus("✍️ Challenge signed");
    } catch {
      setAuthStatus("❌ Signing failed");
    }
  };

  const handleVerify = async () => {
    setLoading(true);
    try {
      const res = await verifyChallenge(challengeId, signature);

      if (res.data.status === "success") {
        setToken(res.data.access_token);
        setAuthStatus("✅ Login successful");
      }
    } catch (err: any) {
      const data = err.response?.data;

      if (!data) {
        setAuthStatus("❌ Network error");
      } else if (data.status === "revoked") {
        setAuthStatus("🚫 Device revoked");
      } else if (data.status === "high_risk") {
        setAuthStatus("⚠️ High risk detected");
      } else {
        setAuthStatus("❌ Verification failed");
      }
    }
    setLoading(false);
  };

  const handleAccess = async () => {
    setLoading(true);
    try {
      const res = await getSecureData(token);
      setResponse(JSON.stringify(res.data));
    } catch (err: any) {
      const data = err.response?.data;

      if (!data) {
        setResponse("❌ Network error");
      } else if (data.status === "revoked") {
        setResponse("🚫 Device revoked");
        setToken("");
      } else {
        setResponse("❌ Access denied");
      }
    }
    setLoading(false);
  };

  const handleRevoke = async () => {
    setLoading(true);
    try {
      const res = await revokeDevice(deviceId);

      if (res.data.status === "success") {
        setRevokeResponse("🚫 Device revoked successfully");

        // reset everything
        setToken("");
        setChallenge("");
        setSignature("");
        setResponse("");
        setAuthStatus("🚫 Access revoked");
      }
    } catch {
      setRevokeResponse("❌ Failed to revoke device");
    }
    setLoading(false);
  };
  const Info = ({ label, value }: { label: string; value: string }) => (
  <div className="bg-gray-800 p-3 rounded-lg break-all">
    <p className="text-gray-400 text-xs">{label}</p>
    <p>{value || "-"}</p>
  </div>
);

  return (
  <div className="min-h-screen bg-gray-950 text-white flex items-center justify-center p-6">
    <div className="w-full max-w-3xl bg-gray-900 rounded-2xl shadow-xl p-6 space-y-6">

      <h1 className="text-2xl font-bold text-center">🔐 Zero Trust System</h1>

      {loading && (
        <div className="text-center text-blue-400 animate-pulse">Loading...</div>
      )}

      <div className="bg-gray-800 p-3 rounded-lg">
        <p className="text-sm text-gray-400">Status</p>
        <p className="font-semibold">{authStatus}</p>
      </div>

      <div className="bg-gray-800 p-3 rounded-lg break-all">
        <p className="text-sm text-gray-400">Public Key</p>
        <p className="text-xs">{publicKey}</p>
      </div>

      {/* Buttons */}
      <div className="grid grid-cols-2 gap-4">

        <button
          disabled={!publicKey}
          onClick={handleRegister}
          className="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-700 px-4 py-2 rounded-lg"
        >
          Register Device
        </button>

        <button
          disabled={!deviceId}
          onClick={handleChallenge}
          className="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-700 px-4 py-2 rounded-lg"
        >
          Get Challenge
        </button>

        <button
          disabled={!challenge}
          onClick={handleSign}
          className="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-700 px-4 py-2 rounded-lg"
        >
          Sign Challenge
        </button>

        <button
          disabled={!signature}
          onClick={handleVerify}
          className="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-700 px-4 py-2 rounded-lg"
        >
          Verify
        </button>

        <button
          disabled={!token}
          onClick={handleAccess}
          className="bg-green-600 hover:bg-green-700 disabled:bg-gray-700 px-4 py-2 rounded-lg col-span-2"
        >
          Access Secure API
        </button>

        <button
          disabled={!deviceId}
          onClick={handleRevoke}
          className="bg-red-600 hover:bg-red-700 disabled:bg-gray-700 px-4 py-2 rounded-lg col-span-2"
        >
          Revoke Device
        </button>

      </div>

      {/* Info */}
      <div className="space-y-3 text-sm">
        <Info label="Device ID" value={deviceId} />
        <Info label="Challenge" value={challenge} />
        <Info label="Signature" value={signature} />
        <Info label="Token" value={token ? token.slice(0, 40) + "..." : ""} />
        <Info label="Response" value={response} />
        <Info label="Revoke Status" value={revokeResponse} />
      </div>

    </div>
  </div>
);
}