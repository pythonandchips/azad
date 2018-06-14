# Azad Spec

## Directory structure

```
root
|- inventory
|- roles
|  |- base_security
|  |  |- tasks
|  |  |  |- main.az
|  |  |  |- pipeline.az
|
|- kibana_server.az

```

## Configuration example

### Playbook configuration

```
server 'kibana' {
  hosts = "tag_kibana_server"
  roles = [
    'base_security',
    'elasticsearch',
    'kibana',
    'logstash',
    'filebeat'
  ]
}
```

### Role configuration

task 'bash', 'whoami' {
  command = 'whoami'
}
