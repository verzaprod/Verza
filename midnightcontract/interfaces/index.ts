export type Bytes32 = string;

export interface RegistryRecord {
  key: Bytes32;
  value: string;
}

export interface AdminEvent {
  admin: string;
  added: boolean;
  timestamp: number;
}

