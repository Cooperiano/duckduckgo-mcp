#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const os = require('os');
const https = require('https');

const VERSION = '2.0.0';
const BASE_URL = `https://github.com/Cooperiano/duckduckgo-mcp/releases/download/v${VERSION}`;

function getPlatform() {
  const platform = os.platform();
  const arch = os.arch();

  if (platform === 'linux' && arch === 'x64') return 'linux-amd64';
  if (platform === 'linux' && arch === 'arm64') return 'linux-arm64';
  if (platform === 'darwin' && arch === 'x64') return 'darwin-amd64';
  if (platform === 'darwin' && arch === 'arm64') return 'darwin-arm64';
  if (platform === 'win32' && arch === 'x64') return 'windows-amd64';

  throw new Error(`Unsupported platform: ${platform}-${arch}`);
}

function downloadFile(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    https.get(url, (response) => {
      if (response.statusCode === 302) {
        // Follow redirect
        https.get(response.headers.location, (redirectResponse) => {
          redirectResponse.pipe(file);
          file.on('finish', () => {
            file.close();
            resolve();
          });
        }).on('error', reject);
      } else {
        response.pipe(file);
        file.on('finish', () => {
          file.close();
          resolve();
        });
      }
    }).on('error', reject);
  });
}

async function install() {
  const platform = getPlatform();
  const binaryName = `duckduckgo-mcp-${platform}`;
  const url = `${BASE_URL}/${binaryName}`;
  const destDir = path.join(__dirname, 'bin');
  const destFile = path.join(destDir, 'duckduckgo-mcp');

  console.log(`Downloading duckduckgo-mcp for ${platform}...`);

  // Create bin directory
  if (!fs.existsSync(destDir)) {
    fs.mkdirSync(destDir, { recursive: true });
  }

  try {
    await downloadFile(url, destFile);
    fs.chmodSync(destFile, '755');
    console.log('Installation complete!');
  } catch (error) {
    console.error('Installation failed:', error.message);
    process.exit(1);
  }
}

install();
