[![Build Status](https://travis-ci.org/mikestaszel/azssh.svg?branch=master)](https://travis-ci.org/mikestaszel/azssh)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# azssh #
Tool to SSH into EC2 instances by name and perform basic operations such as starting, stopping, and rebooting.

This tool gets instance public DNS addresses and "Name" tag values and allows you to list, start, stop, restart, or SSH
into instances.

# Usage #

To use this tool, you must have AWS credentials configured.
[This AWS guide](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html) will help you get set up.

Simply download the binary from the [GitHub release page](https://github.com/mikestaszel/azssh/releases) or clone
and `go build` the source to get started.

### Listing Instances ###

To list instances, run the following:

    azssh ls

### Starting an Instance ###

To start an instance, run the following:

    azssh start instance_name

### Stopping an Instance ###

To stop an instance, run the following:

    azssh stop instance_name

### Restarting an Instance ###

To restart an instance, run the following:

    azssh restart instance_name

### Getting Public DNS of an Instance ###

To print the public DNS address of an instance, run the following:

    azssh dns instance_name

### SSHing into an Instance ###

To SSH into an instance, run the following:

    azssh ssh instance_name

# License #
See the LICENSE file.
