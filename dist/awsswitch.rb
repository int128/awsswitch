class Awsswitch < Formula
  desc "Export the credentials variables to switch a role with MFA"
  homepage "https://github.com/int128/awsswitch"
  url "https://github.com/int128/awsswitch/releases/download/{{ env "VERSION" }}/awsswitch_darwin_amd64.zip"
  version "{{ env "VERSION" }}"
  sha256 "{{ sha256 .darwin_amd64_archive }}"
  def install
    bin.install "awsswitch"
  end
  test do
    system "#{bin}/awsswitch -h"
  end
end
