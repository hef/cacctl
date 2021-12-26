# cacctl

Control your CAC resource from the convenience of your shell

# Install

## Using Homebrew

`brew install hef/tap/cacctl`

## Directly from Github 

Download the latest release from https://github.com/hef/cacctl/releases/latest for your operating system.

# Configuration

| Environment Variable | Description       |
|----------------------|-------------------|
| CAC_USERNAME         | Your CAC Username |
| CAC_PASSWORD         | Your CAC Password |

# Usage

## List Servers

`cacctl list`

### Example Output


```
ID        NAME                          STATUS     IP            CPU RAM SSD PACKAGE
255330813                               Installing 142.47.89.190 1   512 10  CloudPRO v4
255330812 c999963378-cloudpro-329575357 Powered On 142.47.89.189 1   512 10  CloudPRO v4
```

## Build a Server

`cacctl build`

Build a server with 1 CPU, 512MB of ram, and 10GB of storage:

`cactl build --cpu 1 --ram 512 --storage 10`

### Options

| Option     | Type    | Description              |
|------------|---------|--------------------------|
| cpu        | Integer | CPU Count                |
| encryption | Boolean | Encrypt the drive        |
| ha         | Boolean | Enable High Availability |
| os         | String  | Operating System         |
| ram        | Integer | Memory                   |
| storage    | Integer | Disk Space               |

### Valid Operating Systems:

* CentOS 7.9 64Bit
* CentOS 8.3 64bit
* Debian 9.13 64bit
* FreeBSD 12.2 64bit
* Ubuntu 18.04 LTS 64bit

You May have other options available.

## ssh-copy-id

`cacctl ssh-copy-id`
The command will log into all of your instances and copy your ssh public key into authorized_keys
