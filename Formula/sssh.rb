class Sssh < Formula
  desc "SSH into hosts using ~/.ssh/config aliases"
  homepage "https://github.com/dmnkx/homebrew-sssh"
  url "https://github.com/dmnkx/homebrew-sssh/archive/refs/tags/v0.1.4.tar.gz"
  sha256 "fb0e5eaf76a84a52f17f8c9b5a39b32bbe6da2b48aa365d92eb00f926adf1c7c"
  license "MIT"
  head "https://github.com/dmnkx/homebrew-sssh.git", branch: "main"

  livecheck do
    url :stable
    regex(/^v?(\d+(?:\.\d+)+)$/i)
  end

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
  end

  test do
    assert_match "sssh", shell_output("#{bin}/sssh --help")
  end
end
