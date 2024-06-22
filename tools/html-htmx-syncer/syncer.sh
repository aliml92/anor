#!/bin/bash

# Directory paths
HTML_DIR="./html/templates/html/pages"
HTMX_DIR="./html/templates/htmx"

# Check if a filename was provided as an argument
if [ $# -eq 0 ]; then
    echo "Usage: $0 <filename>"
    exit 1
fi

# Get the filename from the command line argument
filename="$1"

# Check if the file exists in the html directory
if [ -e "$HTML_DIR/$filename" ]; then
    # Check if the file exists in the htmx directory
    if [ -e "$HTMX_DIR/$filename" ]; then
        # Copy the content of the html file to the corresponding htmx file
        sed '/{{ define/d' "$HTML_DIR/$filename" | sed '$d' > "$HTMX_DIR/$filename"
        echo "Updated $filename"
    else
        echo "File $filename not found in htmx directory"
    fi
else
    echo "File $filename not found in html directory"
fi
