LDA_VERSION=${LDA_VERSION:-v0.0.9}
OS=${OS:-$(uname | tr '[:upper:]' '[:lower:]')}
ARCH=${ARCH:-$(uname -m)}

# Handle architecture translation if not already set
if [[ "$ARCH" == "x86_64" ]]; then
  ARCH="amd64"
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
  ARCH="arm64"
fi

# Check if wget is available, otherwise fall back to curl
if command -v wget &> /dev/null; then
  downloader="wget -O lda-$OS-$ARCH.tar.gz"
else
  downloader="curl -L -o lda-$OS-$ARCH.tar.gz"
fi

# Download, unzip, and move binary in one go
$downloader https://github.com/devzero-inc/local-developer-analytics/releases/download/$LDA_VERSION/lda-$OS-$ARCH.tar.gz && \
tar -xvf lda-$OS-$ARCH.tar.gz && \
sudo mv lda /usr/local/bin/lda && \
rm lda-$OS-$ARCH.tar.gz
