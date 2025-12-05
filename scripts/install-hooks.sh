#!/bin/bash

# Script to install Git hooks for AlgoShield

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Installing Git hooks for AlgoShield...${NC}"

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo -e "${RED}Error: Not a git repository${NC}"
    exit 1
fi

# Configure git to use custom hooks directory
git config core.hooksPath .githooks

# Make hooks executable
chmod +x .githooks/*

echo -e "${GREEN}✓ Git hooks installed successfully!${NC}"
echo ""
echo -e "${YELLOW}Installed hooks:${NC}"
echo "  • pre-commit:  Runs tests and checks formatting before commit"
echo "  • commit-msg:  Validates commit message format (Conventional Commits)"
echo "  • pre-push:    Runs full test suite before push"
echo ""
echo -e "${YELLOW}To bypass hooks (not recommended):${NC}"
echo "  git commit --no-verify"
echo "  git push --no-verify"
echo ""
echo -e "${GREEN}Happy coding! 🛡️${NC}"

