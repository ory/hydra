# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/vivid64"

  # Hydra
  config.vm.network "forwarded_port", guest: 9000, host: 9000

  # Postgres
  config.vm.network "forwarded_port", guest: 9001, host: 9001

  # Sign in
  config.vm.network "forwarded_port", guest: 3000, host: 3000

  # Sign up
  config.vm.network "forwarded_port", guest: 3001, host: 3001

  config.vm.synced_folder ".", "/home/vagrant/go/src/github.com/ory-am/hydra", owner: "vagrant", group: "vagrant"
  config.vm.provision "shell", path: "./bin/vagrant-provision"
  config.vm.provision "shell", path: "./bin/vagrant-boot", privileged: false, run: "always"
  config.vm.provision "shell", path: "./bin/vagrant-update", privileged: false, run: "always"
end
