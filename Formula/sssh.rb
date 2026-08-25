class Sssh < Formula
  desc "SSH into hosts using ~/.ssh/config aliases"
  homepage "https://github.com/dmnkx/homebrew-sssh"
  url "https://github.com/dmnkx/homebrew-sssh/archive/refs/tags/v0.1.2.tar.gz"
  sha256 "fd151c57782c15afd070b41555475d5963d193cdae253574b501f5b31d41c9dc"
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
