# ipv

CLI tool for finding &amp; connecting to fastest IPVanish VPN server. Written in [Go](https://golang.org/).  
Official IPVanish client is available only for Mac & Windows. This project aims to provide similar capabilities, but only for Linux systems.

## Usage

- ### Connecting to VPN

  - Command - **`sudo ipv connect [flags]`**
  - Connecting to the VPN for the first time generates config file in `$HOME` directory of the current user.
  - After the first run, running the above command reads values from config file for connecting to VPN.
  - To override values from config file, flags can be passed.  
    For more info on available flags, use command - **`sudo ipv connect --help`**

- ### Disconnecting from VPN

  - Command - **`sudo ipv disconnect`**
