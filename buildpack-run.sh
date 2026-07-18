#!/bin/bash
set -e

echo "Installing ntgcalls native library before Go build..."
chmod +x install.sh
./install.sh -n --quiet --skip-summary
