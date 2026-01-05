#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

const PLATFORMS = {
  windows: { GOOS: 'windows', GOARCH: 'amd64', ext: '.exe' },
  linux: { GOOS: 'linux', GOARCH: 'amd64', ext: '' },
  'linux-arm64': { GOOS: 'linux', GOARCH: 'arm64', ext: '' },
  'windows-arm64': { GOOS: 'windows', GOARCH: 'arm64', ext: '.exe' }
};

function buildForPlatform(platform) {
  const config = PLATFORMS[platform];
  if (!config) {
    console.error(`Unknown platform: ${platform}`);
    process.exit(1);
  }

  console.log(`Building for ${platform}...`);
  
  const binDir = path.join(__dirname, '..', 'bin');
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  const outputName = `sentinel-${platform}${config.ext}`;
  const outputPath = path.join(binDir, outputName);

  const env = {
    ...process.env,
    GOOS: config.GOOS,
    GOARCH: config.GOARCH,
    CGO_ENABLED: '0'
  };

  try {
    execSync(
      `go build -ldflags="-s -w" -o "${outputPath}" ./cmd/sentinel`,
      {
        env,
        stdio: 'inherit',
        cwd: path.join(__dirname, '..')
      }
    );
    console.log(`✓ Built ${outputName}`);
  } catch (error) {
    console.error(`✗ Failed to build for ${platform}:`, error.message);
    process.exit(1);
  }
}

function buildAll() {
  console.log('Building for all platforms...');
  Object.keys(PLATFORMS).forEach(platform => {
    buildForPlatform(platform);
  });
  console.log('✓ All builds complete');
}

// Main
const target = process.argv[2] || 'all';

if (target === 'all') {
  buildAll();
} else if (PLATFORMS[target]) {
  buildForPlatform(target);
} else {
  console.error(`Unknown target: ${target}`);
  console.error(`Available targets: all, ${Object.keys(PLATFORMS).join(', ')}`);
  process.exit(1);
}

