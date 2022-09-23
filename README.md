# jy
JSON to YAML to JSON converter tool.

### Build
```
make
```

### How to run
```
jy -json file.json -yaml file.yaml
```
The tool will take whatever file exists and convert it to file which does not
exist, e.g. if `file.json` exists and `file.yaml` does not then it will convert
from json to yaml and vice versa.
