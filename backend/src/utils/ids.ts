import { solidityPackedKeccak256, toBeHex } from 'ethers';

export function genRequestId(userAddr: string, verifierAddr: string, nonce: bigint, timestamp: bigint): string {
  const id = solidityPackedKeccak256(['address','address','uint256','uint256'], [userAddr, verifierAddr, nonce, timestamp]);
  return id; // 0x bytes32
}