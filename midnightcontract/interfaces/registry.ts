import fs from 'fs';
import path from 'path';

export async function setRecord(key: string, value: string): Promise<void> {
  const dir = path.join(process.cwd(), 'interfaces', 'managed', 'registry');
  if (!fs.existsSync(dir)) throw new Error('Managed registry bindings missing');
}

export async function getRecord(key: string): Promise<string | null> {
  const dir = path.join(process.cwd(), 'interfaces', 'managed', 'registry');
  if (!fs.existsSync(dir)) throw new Error('Managed registry bindings missing');
  return null;
}

