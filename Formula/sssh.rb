class Sssh < Formula
  desc "SSH into hosts using ~/.ssh/config aliases"
  homepage "https://github.com/dmnkx/homebrew-sssh"
  version "0.1.5"
  license "MIT"

  livecheck do
    url :homepage
    regex(/^v?(\d+(?:\.\d+)+)$/i)
    strategy :github_latest
  end

  on_macos do
    on_arm do
      url "https://github.com/dmnkx/homebrew-sssh/releases/download/v0.1.5/sssh_0.1.5_darwin_arm64.tar.gz"
      sha256 "e6068acfb2ad11bea22d8506c58ec0eb830a9d67a644bb6b4f15ba644d05bcf0"
    end
    on_intel do
      url "https://github.com/dmnkx/homebrew-sssh/releases/download/v0.1.5/sssh_0.1.5_darwin_amd64.tar.gz"
      sha256 "d3485095ba9cdb5dffbc94a961cc9a3d7d04e8f79738762bf853fcc764b28b2e"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/dmnkx/homebrew-sssh/releases/download/v0.1.5/sssh_0.1.5_linux_arm64.tar.gz"
      sha256 "b658f320a8e3ecd2fa45be5405030c3e7461c698a7956902f094eba1c95dac32"
    end
    on_intel do
      url "https://github.com/dmnkx/homebrew-sssh/releases/download/v0.1.5/sssh_0.1.5_linux_amd64.tar.gz"
      sha256 "ec31d223e400a5c9f5921b05d54def22f0b2305e4b6012e8d5189e026fe6d786"
    end
  end

  def install
    bin.install "sssh"
  end

  test do
    assert_match "sssh", shell_output("#{bin}/sssh --help")
  end
end
