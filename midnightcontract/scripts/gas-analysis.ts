import fs from 'fs';
import path from 'path';

function analyzeCompileLog(logPath: string) {
  const raw = fs.readFileSync(logPath, 'utf-8');
  const rows = Array.from(raw.matchAll(/rows=(\d+)/g)).map(m => Number(m[1]));
  return { circuits: rows.length, totalRows: rows.reduce((a, b) => a + b, 0) };
}

function main() {
  const logFile = path.join(process.cwd(), 'compile.log');
  if (!fs.existsSync(logFile)) {
    console.log('No compile.log found. Run compile with logging to analyze.');
    return;
  }
  const res = analyzeCompileLog(logFile);
  console.log(JSON.stringify(res));
}

main();

