spot-ec2-proxifier
==================

Python script that brings up a spot EC2 instance (of a specified type, region and maximimum bid), tags it, and then launches an SSH tunnel (via plink subprocess) to that instance for use with software such as Proxifier.