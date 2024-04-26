# Hypp

## Tests

```shell
go test $(go list ./... 2>/dev/null | grep -vE 'cmd|jsd')
```

## License

Hypp is published under the AGPL, which can be found [here](./LICENSE).

Hypp is derived from [Hyperapp](https://github.com/jorgebucaran/hyperapp).
Hyperapp is published under the MIT License which is included [here](./hyperapp/LICENSE.md).

Note that Hypp is NOT published under the MIT License.

## Development

Below you'll find the package dependency graph.
Red nodes directly or indirectly import `syscall/js`.

```mermaid
flowchart TD

hypp
subgraph hypp-dir["hypp"]
    subgraph driver-dir["driver"]
        driver-html["html"]
        driver-js["js"]
    end

    subgraph examples-dir["examples"]
        subgraph helloWorld-dir["hello-world"]
            examples-helloWorld-app["app"]
            subgraph helloWorld-cmd["cmd"]
                examples-calculator-cmd-html["html"]
                examples-calculator-cmd-js["js"]
            end
        end
    end

    subgraph tag-dir["tag"]
        subgraph tag-cmd-dir["cmd"]
            tag-cmd-generateTags["generate-tags"]
        end
        tag-html["html"]
        tag-svg["svg"]
    end
end

examples-calculator-cmd-html --> hypp
examples-calculator-cmd-html --> driver-html
examples-calculator-cmd-html --> examples-helloWorld-app

examples-calculator-cmd-js --> driver-js
examples-calculator-cmd-js --> examples-helloWorld-app

tag-html --> hypp
tag-html <-.-> tag-cmd-generateTags
tag-svg --> hypp
tag-svg <-.-> tag-cmd-generateTags

classDef syscallJS fill:#f00;
class examples-calculator-cmd-js,driver-js syscallJS;
```
