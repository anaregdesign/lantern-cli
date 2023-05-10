# lantern-cli

## Install
```shell
git clone https://github.com/anaregdesign/lantern-cli.git
cd lantern-cli
go build
```


## Start
```shell
./lantern-cli --host localhost --port 6380
```

## Example

```
>: put vertex a A
>: get vertex a
{"String_":"a"}
>: put vertex b true
>: get vertex b
{"Bool":true}
>: add edge a b 3.14
>: illuminate neighbor a 3 3 false
{
	"vertices": {
		"a": {
			"Value": {
				"String_": "a"
			}
		},
		"b": {
			"Value": {
				"Bool": true
			}
		}
	},
	"edges": {
		"a": {
			"b": 3.14
		}
	}
}
>: exit
```