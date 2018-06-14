# Overview

This document contains an overview of each package and its responsibilities

## azad

Core classes and operations for azad

## cli

Main command executable

## conn

Conn package manages communication with a server sending commands and running them on the target.

## core

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

## integration

playbooks used for integration testing

## parser

Parse the hcl format playbook and evaluate variables where required

## plugins

Load plugins and retrieve tasks from plugin

## schema

Schema for plugins. Plugins must supply a schema so task defined in playbook can be converted.


