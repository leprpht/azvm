# azvm

A Go CLI for managing Azure Virtual Machines directly from the terminal, built with the Azure SDK and Cobra.

## Features

- Authenticate with Azure CLI
- List Azure Virtual Machines
- Check VM status
- Cross-platform release binaries for macOS, Linux, and Windows

## Installation

Download the latest release for your operating system and architecture from the [GitHub Releases](../../releases) page.

Available release binaries:

- macOS Apple Silicon (`darwin-arm64`)
- macOS Intel (`darwin-amd64`)
- Linux x86_64 (`linux-amd64`)
- Linux ARM64 (`linux-arm64`)
- Windows x86_64 (`windows-amd64`)

Choose the binary that matches your operating system and CPU architecture.

### macOS

Download the appropriate macOS binary from the GitHub Releases page.

For Apple Silicon Macs, use the binary ending in:

```text
darwin-arm64
```

For Intel Macs, use:

```text
darwin-amd64
```

The filename shown below is only an example. **Use the exact filename of the binary you downloaded.**

For example:

```bash
sudo install -m 755 azvm-v0.2.1-darwin-arm64 /usr/local/bin/azvm
```

After installation, you can run:

```bash
azvm login
azvm list
azvm status my-vm
```

### Linux

Download the appropriate Linux binary from the GitHub Releases page.

For x86_64 systems, use the binary ending in:

```text
linux-amd64
```

For ARM64 systems, use:

```text
linux-arm64
```

The filename shown below is only an example. **Use the exact filename of the binary you downloaded.**

For example:

```bash
sudo install -m 755 azvm-v0.2.1-linux-amd64 /usr/local/bin/azvm
```

After installation, you can run:

```bash
azvm login
azvm list
azvm status my-vm
```

### Windows

Download the Windows x86_64 binary from the GitHub Releases page.

The Windows binary ends in:

```text
windows-amd64.exe
```

For example, if the downloaded file is:

```text
azvm-v0.2.1-windows-amd64.exe
```

and place it in a directory on your `PATH`.

One simple option is to create a personal `bin` directory:

```powershell
mkdir "$HOME\bin" -Force
Move-Item .\azvm-v0.2.1-windows-amd64.exe "$HOME\bin\azvm.exe"
```

Then add `$HOME\bin` to your user `PATH`.

After restarting PowerShell, you can run:

```powershell
azvm login
azvm list
azvm status my-vm
```

You can verify that Windows can find the command with:

```powershell
Get-Command azvm
```

## Requirements

- Azure CLI
- An active Azure subscription
- Azure permissions to access and manage Virtual Machines

`azvm` currently uses the Azure CLI for authentication.

Log in with:

```bash
azvm login
```

This runs:

```bash
az login
```

The currently selected Azure subscription is used automatically by `azvm`.

You can verify the active Azure subscription with:

```bash
az account show
```

## Usage

### List VMs

```bash
azvm list
```

Example:

```text
Name: my-vm, Resource Group: my-resource-group, Location: westeurope
```

### Check VM status

```bash
azvm status my-vm
```

Example:

```text
Name: my-vm, Status: VM running
```

## Development

Building from source is only necessary for development or if you want to build `azvm` yourself.

```bash
git clone https://github.com/leprpht/azvm.git
cd azvm
```

### Make targets

Build the application:

```bash
make build
```

Run the application:

```bash
make run
```

Run tests:

```bash
make test
```

Format the code:

```bash
make fmt
```

Update dependencies:

```bash
make tidy
```

Remove build artifacts:

```bash
make clean
```

The development binary is created at:

```text
bin/azvm
```

## Releases

`azvm` uses semantic versioning:

```text
vMAJOR.MINOR.PATCH
```

For example:

```text
v0.2.1
```

Release binaries are automatically built for supported platforms using GitHub Actions when a version tag is pushed.

Supported release targets:

- macOS Apple Silicon
- macOS Intel
- Linux x86_64
- Linux ARM64
- Windows x86_64

## Roadmap

- [x] Azure CLI login
- [x] List VMs
- [x] VM status
- [ ] Start VM
- [ ] Stop VM
- [ ] Restart VM
- [ ] Public IP
- [ ] SSH
- [ ] Remote command execution

## License

azvm is licensed under the MIT License. See [LICENSE](LICENSE) for details.
