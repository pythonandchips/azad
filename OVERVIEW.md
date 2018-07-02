# Overview

This document contains an overview of each package and its responsibilities

## azad

Core classes and operations for azad

## cli

Main command executable

## communicator

Handle loading plugins and returning schema

## conn

Conn package manages communication with a server sending commands and running them on the target.

## communitator/core

Plugin contain basic task required to configure server.

Current task list

- bash: run a basic bash command

Expected task list

- exists: check if a file or directory exists
- user: create a user
- group: create a group
- shell: run a shell command
- copy: copy file on remote or from local to remote
- template: write file to remote based on a template

## communitator/awsinventory

Package to supply server descriptions from amazon.

Current Supported Resources

- ec2

Limitations:

- Only 1 region can be searched currently

## helpers/stringslice

Collection of helper methods for slice of string

## logger

Logging package to output to stdout and stderr

## parser

Parse the hcl format playbook and evaluate variables where required

## plugin

Interfaces and objects to be implemented by a plugin to work with azad

## runner

Core functions for running a playbook

