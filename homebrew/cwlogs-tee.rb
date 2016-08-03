require 'formula'

class CwlogsTee < Formula
  VERSION = '0.1.0'

  homepage 'https://github.com/winebarrel/cwlogs-tee'
  url "https://github.com/winebarrel/cwlogs-tee/releases/download/v#{VERSION}/cwlogs-tee-v#{VERSION}-darwin-amd64.gz"
  sha256 'cf318f404f9f62b0fcf3a035433121712e3d14da6545da603521c3e8e4a8260a'
  version VERSION
  head 'https://github.com/winebarrel/cwlogs-tee.git', :branch => 'master'

  def install
    system "mv cwlogs-tee-v#{VERSION}-darwin-amd64 cwlogs-tee"
    bin.install 'cwlogs-tee'
  end
end
