package html

// The tags are taken from https://developer.mozilla.org/en-US/docs/Web/HTML/Element

//go:generate bash -c "go run ../cmd/generate-tags/main.go html html base head link meta style title body address article aside footer header h1 h2 h3 h4 h5 h6 main nav section body blockquote dd div dl dt figcaption figure hr li ol p pre ul a re abbr b bdi bdo br cite code data dfn em i kbd mark q rp rt ruby s samp small span strong sub sup time u var wbr area audio img map track video embed iframe object param picture portal source svg svg canvas noscript script del ins caption col colgroup table tbody td tfoot th thead tr button datalist fieldset form input label legend meter optgroup option output progress select textarea details dialog menu summary slot template > generated.go"
