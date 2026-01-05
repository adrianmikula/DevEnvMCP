#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const os = require('os');

// Determine the correct binary based on platform
const platform = os.platform();
const arch = os.arch();

let binaryName;
if (platform === 'win32') {
  if (arch === 'arm64') {
    binaryName = 'sentinel-windows-arm64.exe';
  } else {
    binaryName = 'sentinel-windows.exe';
  }
} else if (platform === 'linux') {
  if (arch === 'arm64') {
    binaryName = 'sentinel-linux-arm64';
  } else {
    binaryName = 'sentinel-linux';
  }
} else {
  console.warn(`Unsupported platform: ${platform} ${arch}`);
  console.warn('This package is designed for Windows and Linux only.');
  process.exit(0);
}

const binDir = path.join(__dirname, '..', 'bin');
const binaryPath = path.join(binDir, binaryName);

if (!fs.existsSync(binaryPath)) {
  console.error(`Binary not found: ${binaryPath}`);
  console.error('Please run "npm run build" first.');
  process.exit(1);
}

// Make binary executable on Linux
if (platform !== 'win32') {
  try {
    fs.chmodSync(binaryPath, 0o755);
  } catch (error) {
    console.warn(`Warning: Could not make binary executable: ${error.message}`);
  }
}

console.log(`✓ Postinstall complete. Binary: ${binaryName}`);

