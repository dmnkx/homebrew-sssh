#!/usr/bin/env bash
set -euo pipefail

tag="${1:?usage: update-formula.sh v1.2.3}"
url="https://github.com/dmnkx/homebrew-sssh/archive/refs/tags/${tag}.tar.gz"
sha="$(curl -fsSL "$url" | shasum -a 256 | awk '{print $1}')"
mkdir -p Formula

cat > Formula/sssh.rb <<EOF
class Sssh < Formula
  desc "SSH into hosts using ~/.ssh/config aliases"
  homepage "https://github.com/dmnkx/homebrew-sssh"
  url "${url}"
  sha256 "${sha}"
  license "MIT"
  head "https://github.com/dmnkx/homebrew-sssh.git", branch: "main"

  livecheck do
    url :stable
    regex(/^v?(\\d+(?:\\.\\d+)+)\$/i)
  end

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
  end

  test do
    assert_match "sssh", shell_output("#{bin}/sssh --help")
  end
end
EOF
