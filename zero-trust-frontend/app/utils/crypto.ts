import nacl from "tweetnacl";
import * as util from "tweetnacl-util";

export const generateKeyPair = () => {
  const keyPair = nacl.sign.keyPair();

  return {
    publicKey: util.encodeBase64(keyPair.publicKey),
    privateKey: util.encodeBase64(keyPair.secretKey),
  };
};

export const signChallenge = (
  challengeBase64: string,
  privateKeyBase64: string
) => {
  const challenge = util.decodeBase64(challengeBase64);
  const privateKey = util.decodeBase64(privateKeyBase64);

  const signature = nacl.sign.detached(challenge, privateKey);
  return util.encodeBase64(signature);
};