# go-manidator

Manidator is a utility to display the last line of output from one or more "dators".

A "dator" (from "updater") is simply something that can provide live line-based
output by implementing the `manidator.Dator` interface.

## Why does manidator exist?

I wanted to provide output from multiple commands at the same time, but I didn't
want the output of each command to get mixed up with the output of the others.

## Where does the name come from?

The name `manidator` is a sort of portmanteau of "many" and "updater".
