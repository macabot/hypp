// Run the following in the browser.
function getTagOwnProperties(tags) {
    const hierarchy = {}
    const ownProps = {}

    tags.forEach(tag => {
        let o = document.createElement(tag)
        let prevName = o.constructor.name
        hierarchy[tag] = prevName
        let props = Object.getOwnPropertyNames(o)
        if (props.length) {
            ownProps[prevName] = props
        }
        while (o) {
            o = Object.getPrototypeOf(o)
            if (!o) {
                break
            }
            const name = o.constructor.name
            if (name !== prevName) {
                hierarchy[prevName] = name
            }
            if (!ownProps[name]) {
                const props = Object.getOwnPropertyNames(o)
                if (props.length) {
                    ownProps[name] = props
                }
            }
            prevName = name
        }
    })

    return {hierarchy, ownProps}
}

const tags = [] // Fill in HTML tags and SVG tags.
console.log(JSON.stringify(getTagOwnProperties(tags), null, 2))
// Copy output and paste in tag.json.
