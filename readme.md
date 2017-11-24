# Logc
A log collector(like logstash, fluentd) written by golang.
Now only support read from files and output to elasticSearch.

# Install

```
sudo wget -O /usr/local/bin/logc 'https://github.com/lovego/logc/releases/download/170706/logc'
sudo chmod +x /usr/local/bin/logc
```

# Usage
```
logc <yaml_config_file>
```
See <a href="blob/master/testdata/logc.yml">logc.yml</a> for full config format.


