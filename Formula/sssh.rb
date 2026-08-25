class Sssh < Formula
  desc "SSH into hosts using ~/.ssh/config aliases"
  homepage "https://github.com/dmnkx/homebrew-sssh"
  url "https://github.com/dmnkx/homebrew-sssh/archive/refs/tags/v0.1.3.tar.gz"
  sha256 "1d29b2782d342ab67244c7407dcea36a53184f8228c9f91d48e337461dd1617e"
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
