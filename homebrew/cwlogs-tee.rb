require 'formula'

class CwlogsTee < Formula
  VERSION = '0.1.4'

  homepage 'https://github.com/winebarrel/cwlogs-tee'
  url "https://github.com/winebarrel/cwlogs-tee/releases/download/v#{VERSION}/cwlogs-tee-v#{VERSION}-darwin-amd64.gz"
  sha256 '0c9441b883724327c79072fe19cc88e69591e2f73c3210249335fd5c14bb5eb8'
  version VERSION
  head 'https://github.com/winebarrel/cwlogs-tee.git', :branch => 'master'

  def install
    system "mv cwlogs-tee-v#{VERSION}-darwin-amd64 cwlogs-tee"
    bin.install 'cwlogs-tee'
  end
end
