#!/usr/bin/env node

const { execSync } = require('child_process');
const path = require('path');
const fs = require('fs');

const binaryPath = path.join(__dirname, 'duckduckgo-mcp');

// Check if binary exists
if (!fs.existsSync(binaryPath)) {
  console.error('Binary not found. Please reinstall the package.');
  process.exit(1);
}

// Run the binary
try {
  execSync(binaryPath, { stdio: 'inherit' });
} catch (error) {
  process.exit(error.status);
}
