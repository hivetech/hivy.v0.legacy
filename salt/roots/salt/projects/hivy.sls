dependencies:
  pkg.installed:
    - reload_modules: true
    - pkgs:
      - git-core
      - rake
      - bzr
      - mercurial
  gem.installed:
    - name: awesome_print
  file.directory:
    - name: /root/goworkspace/src/github.com/hivetech
    - user: root
    - mode: 755
    - makedirs: True
  cmd.run: 
    - name: sudo cp -r /app /root/goworkspace/src/github.com/hivetech/hivy 

install:
  cmd.run:
    - env: GOPATH=/root/goworkspace
    - name: | 
      rake install:go
      rake install
    - cwd: /app
    - user: root
    - require:
      - pkg: dependencies

startup:
  cmd.run:
    - name: |
      PATH=$PATH:/root/goworkspace/bin GOPATH=/root/goworkspace rake app:init
      cd hivy && PATH=$PATH:/root/goworkspace/bin GOPATH=/root/goworkspace forego start &
    - cwd: /root/goworkspace/src/github.com/hivetech/hivy
    - require:
      - cmd: install
