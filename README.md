# Hypp

## Tests

```shell
$ go test . ./driver/html/... ./examples/.../app ./examples/.../html ./tag/...
```

## License

Hypp is published under the AGPL, which can be found [here](./LICENSE).

Hypp is derived from [Hyperapp](https://github.com/jorgebucaran/hyperapp).
Hyperapp is published under the MIT License which is included [here](./hyperapp/LICENSE.md).

Note that Hypp is NOT published under the MIT License.

## Development

Red nodes directly or indirectly import `syscall/js`.

Current package dependency graph:

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
                examples-calculater-cmd-html["html"]
                examples-calculater-cmd-js["js"]
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

examples-calculater-cmd-html --> hypp
examples-calculater-cmd-html --> driver-html
examples-calculater-cmd-html --> examples-helloWorld-app

examples-calculater-cmd-js --> driver-js
examples-calculater-cmd-js --> examples-helloWorld-app

tag-html --> hypp
tag-html <-.-> tag-cmd-generateTags
tag-svg --> hypp
tag-svg <-.-> tag-cmd-generateTags

classDef syscallJS fill:#f00;
class examples-calculater-cmd-js,driver-js syscallJS;
```

Possible new package dependency graph:

```mermaid
flowchart TD

hypp
window
util
html
svg
js
jsd

subgraph cmd-dir["cmd"]
    cmd-generateTags["generate-tags"]
end

subgraph examples-dir["examples"]
    subgraph examples-hello-world-dir["hello-world"]
        examples-hello-world-app["app"]
        subgraph examples-hello-world-cmd-dir["cmd"]
            examples-hello-world-cmd-js["js"]
            examples-hello-world-cmd-html["html"]
        end
    end
end

examples-hello-world-app --> window
examples-hello-world-cmd-js --> hypp
examples-hello-world-cmd-js --> examples-hello-world-app
examples-hello-world-cmd-js --> jsd
examples-hello-world-cmd-html --> hypp
examples-hello-world-cmd-html --> examples-hello-world-app

html <-.-> cmd-generateTags
svg <-.-> cmd-generateTags
html --> hypp
svg --> hypp
jsd --> js
hypp --> util
hypp --> window
window --> util
window --> js

classDef syscallJS fill:#f00;
class jsd,examples-hello-world-cmd-js syscallJS;
```
