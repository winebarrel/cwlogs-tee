require 'formula'

class CwlogsTee < Formula
  VERSION = '0.1.1'

  homepage 'https://github.com/winebarrel/cwlogs-tee'
  url "https://github.com/winebarrel/cwlogs-tee/releases/download/v#{VERSION}/cwlogs-tee-v#{VERSION}-darwin-amd64.gz"
  sha256 '03c1448b740fdf92833ff1842ef15dccc111e2eeb74642e12188380b08226925'
  version VERSION
  head 'https://github.com/winebarrel/cwlogs-tee.git', :branch => 'master'

  def install
    system "mv cwlogs-tee-v#{VERSION}-darwin-amd64 cwlogs-tee"
    bin.install 'cwlogs-tee'
  end
end
