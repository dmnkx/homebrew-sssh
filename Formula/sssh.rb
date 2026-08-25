class Sssh < Formula
  desc "SSH into hosts using ~/.ssh/config aliases"
  homepage "https://github.com/dmnkx/homebrew-sssh"
  version "0.1.4"
  license "MIT"

  livecheck do
    url :homepage
    regex(/^v?(\d+(?:\.\d+)+)$/i)
    strategy :github_latest
  end

  on_macos do
    on_arm do
      url "https://github.com/dmnkx/homebrew-sssh/releases/download/v0.1.4/sssh_0.1.4_darwin_arm64.tar.gz"
      sha256 "c003cbff0115756ba97a79322a0d577a2d0014a423018e740aa417df01fd69a0"
    end
    on_intel do
      url "https://github.com/dmnkx/homebrew-sssh/releases/download/v0.1.4/sssh_0.1.4_darwin_amd64.tar.gz"
      sha256 "7cd7f9b133a62c24135677ae1b3f738102327e85e992c7d995db560dea0a161a"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/dmnkx/homebrew-sssh/releases/download/v0.1.4/sssh_0.1.4_linux_arm64.tar.gz"
      sha256 "793870bec47e2ba86d8561caa501e0101ef9bfecd572ef790864e8212d204732"
    end
    on_intel do
      url "https://github.com/dmnkx/homebrew-sssh/releases/download/v0.1.4/sssh_0.1.4_linux_amd64.tar.gz"
      sha256 "deb61c4f017530b0654df5c4650e15b907c5b9b922ebf98a71fc1677ec8349cf"
    end
  end

  def install
    bin.install "sssh"
  end

  test do
    assert_match "sssh", shell_output("#{bin}/sssh --help")
  end
end
