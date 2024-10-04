name: sshx
prefix: ""

topology:
  nodes:
    sshx:
      kind: nokia_srlinux
      image: ghcr.io/nokia/srlinux:24.7.2
      exec:
        - touch /tmp/.ndk-dev-mode
        {{- if ne (env.Getenv "NDK_DEBUG") "" }}
        - /debug/prepare-debug.sh
        {{- end }}
      binds:
        - ../build:/tmp/build # mount app binary
        - ../sshx.yml:/tmp/sshx.yml # agent config file to appmgr directory
        - ../yang:/opt/sshx/yang # yang modules
        - ../logs/srl:/var/log/srlinux # expose srlinux logs
        - ../logs/sshx/:/var/log/sshx # expose greeter log file
        - ../bin/sshx-0.2.5:/opt/sshx/sshx-bin # sshx binary
        {{- if ne (env.Getenv "NDK_DEBUG") "" }}
        - ../debug/:/debug/
        {{- end }}