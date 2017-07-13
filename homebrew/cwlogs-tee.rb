require 'formula'

class CwlogsTee < Formula
  VERSION = '0.1.4'

  homepage 'https://github.com/winebarrel/cwlogs-tee'
  url "https://github.com/winebarrel/cwlogs-tee/releases/download/v#{VERSION}/cwlogs-tee-v#{VERSION}-darwin-amd64.gz"
  sha256 '0cad6dee396540d0fd2fd97cdbcaf160b2df9bacf333392de22e26a24ae0d4d9'
  version VERSION
  head 'https://github.com/winebarrel/cwlogs-tee.git', :branch => 'master'

  def install
    system "mv cwlogs-tee-v#{VERSION}-darwin-amd64 cwlogs-tee"
    bin.install 'cwlogs-tee'
  end
end
