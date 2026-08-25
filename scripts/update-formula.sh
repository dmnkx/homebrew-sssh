#!/usr/bin/env bash
set -euo pipefail

tag="${1:?usage: update-formula.sh v1.2.3}"
ver="${tag#v}"
base="https://github.com/dmnkx/homebrew-sssh/releases/download/${tag}"
sums="$(curl -fsSL "${base}/checksums.txt")"

sha_for() {
  local file="$1"
  echo "${sums}" | awk -v f="${file}" '$2 == f { print $1; exit }'
}

sha_darwin_arm="$(sha_for "sssh_${ver}_darwin_arm64.tar.gz")"
sha_darwin_amd="$(sha_for "sssh_${ver}_darwin_amd64.tar.gz")"
sha_linux_arm="$(sha_for "sssh_${ver}_linux_arm64.tar.gz")"
sha_linux_amd="$(sha_for "sssh_${ver}_linux_amd64.tar.gz")"

if [[ -z "${sha_darwin_arm}" || -z "${sha_darwin_amd}" || -z "${sha_linux_arm}" || -z "${sha_linux_amd}" ]]; then
  echo "missing checksums for ${tag}" >&2
  echo "${sums}" >&2
  exit 1
fi

mkdir -p Formula
cat > Formula/sssh.rb <<EOF
class Sssh < Formula
  desc "SSH into hosts using ~/.ssh/config aliases"
  homepage "https://github.com/dmnkx/homebrew-sssh"
  version "${ver}"
  license "MIT"

  livecheck do
    url :homepage
    regex(/^v?(\\d+(?:\\.\\d+)+)\$/i)
    strategy :github_latest
  end

  on_macos do
    on_arm do
      url "${base}/sssh_${ver}_darwin_arm64.tar.gz"
      sha256 "${sha_darwin_arm}"
    end
    on_intel do
      url "${base}/sssh_${ver}_darwin_amd64.tar.gz"
      sha256 "${sha_darwin_amd}"
    end
  end

  on_linux do
    on_arm do
      url "${base}/sssh_${ver}_linux_arm64.tar.gz"
      sha256 "${sha_linux_arm}"
    end
    on_intel do
      url "${base}/sssh_${ver}_linux_amd64.tar.gz"
      sha256 "${sha_linux_amd}"
    end
  end

  def install
    bin.install "sssh"
  end

  test do
    assert_match "sssh", shell_output("#{bin}/sssh --help")
  end
end
EOF
