ec2-spot-instance-launcher
==================
This is a Python script that brings up a spot EC2 instance (of a specified type, region and maximimum bid), and tags it. It also allows to list and stop active instances.

Prerequisites
---------
* Python 2.6 (ish)
* [Boto](https://github.com/boto/boto) library

Configuration
---------
A template configuration file is provided in `ec2-spot-instance-launcher.cfg.defaults` (rename it to `ec2-spot-instance-launcher.cfg`).

These are some important parameters:

- `ami` Unfortunately you will need to manually find the AMI to use in your region
- `key_pair` You will need a named keypair on EC2 for instance creation
- `security_group` Ditto
- `tag` A tag that will be used by the script to identify the instance's purpose
- `user_data` User data be passed by the script to instance
- `user_data_file` User data file would be used if user_data is not set. Script would read the file content, and pass the same to instance as user-data


Some thoughts
--------------

IAM user needs to have access to a number of EC2 API calls, I believe they are:
* DescribeInstances
* DescribeSpotPriceHistoryRequest
* DescribeSpotInstanceRequests
* RequestSpotInstances

The script is written in a fairly platform agnostic way but it's only been tested on Mac OSX.

Expected behavior
-------------
First connect  should go something like:
<pre>python main.py
Spot price is 0.02... below maximum bid, continuing
Spot request created, status: open
Waiting for instance to be provisioned (usually takes 1m to be reviewed, another 2m to be fulfilled) ...  . . . Instance is active.
Tagging instance... Done.
Waiting for server to come up ... . . . . Server is up!
</pre>

Subsequent connects:
<pre>python main.py
Instance exists already, we will not be provisioning another one
Waiting for server to come up ... Server is up!
</pre>

Terminating (and detagging) the instance:
<pre>python main.py stop
Terminating i-0a1b1930 ... done.
</pre>

Credits
-------------
Most code comes from [spot-ec2-proxifier](https://github.com/alexzorin/spot-ec2-proxifier) project
