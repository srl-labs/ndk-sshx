sshx:
  path: /usr/local/bin
  launch-command: {{if ne (env.Getenv "NDK_DEBUG") "" }}{{ "/debug/dlv --listen=:7000"}}{{ if ne (env.Getenv "NDK_DEBUG") "" }} {{ "--continue --accept-multiclient" }}{{ end }} {{ "--headless=true --log=true --api-version=2 exec"}} {{ end }}sshx
  version-command: sshx --version
  failure-action: wait=10
  config-delivery-format: json
  yang-modules:
    names:
      - sshx
    source-directories:
      - /opt/sshx/yang
