require 'formula'

class CwlogsTee < Formula
  VERSION = '0.1.2'

  homepage 'https://github.com/winebarrel/cwlogs-tee'
  url "https://github.com/winebarrel/cwlogs-tee/releases/download/v#{VERSION}/cwlogs-tee-v#{VERSION}-darwin-amd64.gz"
  sha256 'e4c6943a3ee8b49c9a4f40b287ac13dac267215611658d4aabe0a6cfabc0fe39'
  version VERSION
  head 'https://github.com/winebarrel/cwlogs-tee.git', :branch => 'master'

  def install
    system "mv cwlogs-tee-v#{VERSION}-darwin-amd64 cwlogs-tee"
    bin.install 'cwlogs-tee'
  end
end
