# cacctl

Control your CAC resource from the convenience of your shell

# Install

## Using homebrew

`brew install hef/tap/cacctl`

# Configuration

| Environment Variable  | Description        |
| --------------------- | ------------------ |
| CAC_USERNAME          | Your CAC Username  |
| CAC_PASSWORD          | Your CAC Password  |

# Usage

## List Servers

`cacctl list`

### Example Output


```
ID        NAME                          IP            CPU RAM SSD PACKAGE
255330777 Installing                    142.47.89.170 1   512 10  
255330776 c999963378-cloudpro-728834606 142.47.89.169 1   512 10  CloudPRO v4
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

You May have other options available.
