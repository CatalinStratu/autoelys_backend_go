#!/bin/bash

echo "=== Brand Logo Download Script ==="

# API endpoint
API_URL="http://localhost:8080/api/brands"

# Directory to save logos
LOGO_DIR="public/logos/brands"

# Create directory if it doesn't exist
mkdir -p "$LOGO_DIR"

# Fetch brands from API
echo "Fetching brands from API..."
BRANDS_JSON=$(curl -s -X 'GET' "$API_URL" -H 'accept: application/json')

if [ $? -ne 0 ]; then
    echo "❌ Failed to fetch brands from API"
    exit 1
fi

# Check if we got data
if [ -z "$BRANDS_JSON" ]; then
    echo "❌ No data received from API"
    exit 1
fi

echo "✅ Successfully fetched brands"

# Parse JSON and download logos
echo ""
echo "Downloading brand logos..."

# Extract brand data and download logos
echo "$BRANDS_JSON" | jq -r '.data[] | "\(.id)|\(.name)|\(.logo)"' | while IFS='|' read -r id name logo; do
    if [ "$logo" != "null" ] && [ -n "$logo" ]; then
        # Clean brand name for filename (remove spaces, convert to lowercase)
        filename=$(echo "$name" | tr '[:upper:]' '[:lower:]' | tr ' ' '-')

        # Get file extension from URL
        extension="${logo##*.}"
        # Handle query parameters in URL
        extension=$(echo "$extension" | cut -d'?' -f1)

        # If extension is too long or invalid, default to jpg
        if [ ${#extension} -gt 4 ]; then
            extension="jpg"
        fi

        filepath="$LOGO_DIR/${filename}.${extension}"

        # Download logo
        echo "  [$id] Downloading $name logo..."
        curl -s -o "$filepath" "$logo"

        if [ $? -eq 0 ]; then
            echo "  ✅ Saved: $filepath"
        else
            echo "  ❌ Failed to download: $name"
        fi
    else
        echo "  ⚠️  No logo for: $name"
    fi
done

echo ""
echo "=== Download Complete ==="

# Count downloaded files
file_count=$(ls -1 "$LOGO_DIR" | wc -l)
echo "Total logos downloaded: $file_count"

# Create a mapping file
echo ""
echo "Creating brand-logo mapping file..."

echo "$BRANDS_JSON" | jq -r '.data[] | "\(.id),\(.name),\(.logo)"' > "${LOGO_DIR}/brand-logo-mapping.csv"

if [ $? -eq 0 ]; then
    echo "✅ Mapping file created: ${LOGO_DIR}/brand-logo-mapping.csv"
else
    echo "❌ Failed to create mapping file"
fi

echo ""
echo "All done! Logos saved in: $LOGO_DIR"
