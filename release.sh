#!/bin/bash

rm -rf example/ssh/dist

pushd example/ssh/
	echo "build bin ${BUILD_VERSION} ${BUILD_DATE} ${COMMIT_SHA1}"
	gox -osarch="darwin/amd64 linux/amd64 windows/amd64" \
        -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}"
popd

pushd example/ssh/dist
sha256sum ssh_darwin_amd64 > ssh_darwin_amd64.sha256sum
sha256sum ssh_linux_amd64 > ssh_linux_amd64.sha256sum
sha256sum ssh_windows_amd64.exe > ssh_windows_amd64.sha256sum
popd

version=$(cat version.txt)
macsha=$(cat example/ssh/dist/ssh_darwin_amd64.sha256sum | awk '{print $1}')
linuxsha=$(cat example/ssh/dist/ssh_linux_amd64.sha256sum | awk '{print $1}')

cat > ysssh.rb <<EOF
class Ysssh < Formula
    desc "Devops tools 运维工具 ysssh"
    homepage "https://github.com/ysicing/ssh"
    version "${version}"
    bottle :unneeded

    if OS.mac?
      url "https://github.com/ysicing/ssh/releases/download/#{version}/ssh_darwin_amd64"
      sha256 "${macsha}"
    elsif OS.linux?
      if Hardware::CPU.intel?
        url "https://github.com/ysicing/ssh/releases/download/#{version}/ssh_linux_amd64"
        sha256 "${macsha}"
      end
    end

    def install
      bin.install "ssh_darwin_amd64" => "ysssh"
    end
  end
EOF

make release