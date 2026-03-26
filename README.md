kubernetes-live-images
======================

Storing docker images that aren't needed costs money unnecessarily. Not storing
images you still need can cause outages. ECR (and other equivalents) can be
configured to delete images not pulled for some period of time, but if a pod is
left running for longer than that then its image may be deleted, leading to
breakage if e.g. the node hosting the pod then goes down.

This tool examines all of the pods on a cluster, and pulls the metadata for
each image used. If run on a regular schedule, this should prevent images from
being cleaned up while they are still in use.
