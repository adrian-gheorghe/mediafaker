#!/bin/bash

# Get final URL after curl is redirected
RELEASE_URL=$(curl -Ls -o /dev/null -w %{url_effective} https://github.com/adrian-gheorghe/mediafaker/releases/latest)
# Extract tag after the last forward slash
TAG="${RELEASE_URL##*/}"

# Check if mediafaker is currently installed
LOCAL_PATH=$(which mediafaker)
 if [ -x "$LOCAL_PATH" ]; then
    echo "Mediafaker is already Installed"
    echo "Try running $(mediafaker --help)"
    exit 0
fi

echo "Attempting to download mediafaker v${TAG}"

if [[ "$OSTYPE" == "linux-gnu" ]]; then
    BINARY_PATH="mediafaker-linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    BINARY_PATH="mediafaker-darwin"
else
    BINARY_PATH="mediafaker"
fi

curl -L "https://github.com/adrian-gheorghe/mediafaker/releases/download/$TAG/$BINARY_PATH" --output $BINARY_PATH
chmod +x $BINARY_PATH
mv $BINARY_PATH /usr/local/bin/mediafaker

echo "Mediafaker installed successfully!"
mediafaker --help