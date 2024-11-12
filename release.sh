#!/bin/bash

BINARY_NAME=dns-manager
VERSION=$(git describe --tags $(git rev-list --tags --max-count=1))
VERSION_WITHOUT_V=$(echo $VERSION | cut -c 2-)
SHA256_DARWIN_ARM64=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-darwin-arm64.tar.gz | awk '{print $1}')
SHA256_DARWIN_AMD64=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-darwin-amd64.tar.gz | awk '{print $1}')
SHA256_LINUX_ARM64=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-linux-arm64.tar.gz | awk '{print $1}')
SHA256_LINUX_ARM=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-linux-arm.tar.gz | awk '{print $1}')
SHA256_LINUX_AMD64=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-linux-amd64.tar.gz | awk '{print $1}')
SHA256_LINUX_386=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-linux-386.tar.gz | awk '{print $1}')
SHA256_WIN_AMD64=$(shasum -a 256 builds/$BINARY_NAME-$VERSION-windows-amd64 | awk '{print $1}')
WIN_AMD64="$BINARY_NAME-$VERSION-windows-amd64"

# Generate formula from template with replacements
cat <<EOF > dns-manager.rb
class Dnsmanager < Formula
  desc "dns-manager!"
  homepage "https://www.kuepper.nrw"
  
  version "${VERSION_WITHOUT_V}"
  license "MIT"

  test do
  end

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ruedigerp/cloudflare-dns-manager-homebrew/releases/download/${VERSION}/dns-manager-${VERSION}-darwin-arm64.tar.gz"
      sha256 "$SHA256_DARWIN_ARM64"
    elsif Hardware::CPU.intel?
      url "https://github.com/ruedigerp/cloudflare-dns-manager-homebrew/releases/download/${VERSION}/dns-manager-${VERSION}-darwin-amd64.tar.gz"
      sha256 "$SHA256_DARWIN_AMD64"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/ruedigerp/cloudflare-dns-manager-homebrew/releases/download/${VERSION}/dns-manager-${VERSION}-linux-amd64.tar.gz"
        sha256 "$SHA256_LINUX_AMD64"
      else
        url "https://github.com/ruedigerp/cloudflare-dns-manager-homebrew/releases/download/${VERSION}/dns-manager-${VERSION}-linux-386.tar.gz"
        sha256 "$SHA256_LINUX_386"
      end
    elsif Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/ruedigerp/cloudflare-dns-manager-homebrew/releases/download/${VERSION}/dns-manager-${VERSION}-linux-arm64.tar.gz"
        sha256 "$SHA256_LINUX_ARM64"
      else
        url "https://github.com/ruedigerp/cloudflare-dns-manager-homebrew/releases/download/${VERSION}/dns-manager-${VERSION}-linux-arm.tar.gz"
        sha256 "$SHA256_LINUX_ARM"
      end
    end
  end
  
  def install
  if OS.mac?
    if Hardware::CPU.arm?
      # Installation steps for macOS ARM64
      bin.install "dns-manager-$VERSION-darwin-arm64" => "dns-manager"
    elsif Hardware::CPU.intel?
      # Installation steps for macOS AMD64
      bin.install "dns-manager-$VERSION-darwin-amd64" => "dns-manager"
    end
  elsif OS.linux?
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        # Installation steps for Linux AMD64
        bin.install "dns-manager-$VERSION-linux-amd64" => "dns-manager"
      else
        # Installation steps for Linux 386
        bin.install "dns-manager-$VERSION-linux-386" => "dns-manager"
      end
    elsif Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        # Installation steps for Linux ARM64
        bin.install "dns-manager-$VERSION-linux-arm64" => "dns-manager"
      else
        # Installation steps for Linux ARM
        bin.install "dns-manager-$VERSION-linux-arm" => "dns-manager"
      end
    end
  end
end
end
EOF


cat <<EOF > dns-manager.json
{
    "version": "$VERSION_WITHOUT_V",
    "license": "UNKNOWN",
    "homepage": "https://kuepper.nrw",
    "bin": "dns-manager",
    "pre_install": "Rename-Item \"\$dir\\\\$WIN_AMD64\" dns-manager",
    "description": "Dnsmanager!",
    "architecture": {
        "64bit": {
            "url": "https://github.com/ruedigerp/cloudflare-dns-manager/releases/download/$VERSION/dns-manager-$VERSION-windows-amd64",
            "hash": "$SHA256_WIN_AMD64"
        }
    }
}
EOF