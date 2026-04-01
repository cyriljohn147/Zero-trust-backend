import axios from "axios";

const API = axios.create({
  baseURL: "http://localhost:8080",
});

export const registerDevice = (publicKey: string) =>
  API.post("/devices/register", { user_id: 1, public_key: publicKey });

export const getChallenge = (deviceId: string) =>
  API.post("/auth/challenge", { device_id: deviceId });

export const signChallenge = (challenge: string) =>
  API.post("/auth/sign", { challenge });

export const verifyChallenge = (challengeId: string, signature: string) =>
  API.post("/auth/verify", {
    challenge_id: challengeId,
    signature,
  });

export const getSecureData = (token: string) =>
  API.get("/api/secure-data", {
    headers: { Authorization: `Bearer ${token}` },
  });

export const revokeDevice = (deviceId: string) =>
  API.post("/devices/revoke", { device_id: deviceId });
