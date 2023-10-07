class Ogpk < Formula
  desc "CLI tool to fetch OpenGraph data from a URL"
  homepage "https://github.com/almonk/ogpk"
  url "https://github.com/almonk/ogpk/releases/download/0.1.0/ogpk-0.1.0-darwin-amd64"
  sha256 "6d81c9928f845c3fc257c07600303c8ed5a4ba80d9a90c7c3d1de5e672674281"
  version "0.1"

  def install
    bin.install "ogpk"
  end
end
