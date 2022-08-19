package htmlx

import (
	"strings"
	"testing"

	"github.com/starudream/go-lib/testx"
)

var raw = `
<!doctype html>
<html lang="zh">

<head>
    <meta charset="UTF-8">
    <meta content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0" name="viewport">
    <meta content="ie=edge" http-equiv="X-UA-Compatible">
    <title>
TEST
</title>
    <link crossorigin="anonymous" href="https://cdn.jsdelivr.net/npm/semantic-ui-css@2.4.1/semantic.min.css" integrity="sha256-9mbkOfVho3ZPXfM7W8sV2SndrGDuh7wuyLjtsWeTI1Q=" rel="stylesheet">
</head>

<body>
<div style="height: 50px;"></div>
<div class="ui container">
    <div class="ui centered grid">
        <table class="ui large celled structured table">
            <thead>
            <tr>
                <th style="width: 400px;">Category</th>
                <th style="width: 400px;">Item</th>
                <th>Desc</th>
            </tr>
            </thead>
            <tbody id="list">
            </tbody>
        </table>
    </div>
</div>
<div style="height: 200px;"></div>
</body>
<script crossorigin="anonymous" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" src="https://cdn.jsdelivr.net/npm/jquery@3.6.0/dist/jquery.min.js"></script>
<script crossorigin="anonymous" integrity="sha256-t8GepnyPmw9t+foMh3mKNvcorqNHamSKtKRxxpUEgFI=" src="https://cdn.jsdelivr.net/npm/semantic-ui-css@2.4.1/semantic.min.js"></script>
<script>
    const list = $("#list")
</script>

</html>
`

var root *Node

func init() {
	node, err := Parse(strings.NewReader(raw))
	if err != nil {
		panic(err)
	}
	root = node
}

func TestNodeSearch(t *testing.T) {
	title := NodeSearch(root, func(node *Node) bool {
		return node.Type == ElementNode && node.Data == "title"
	})
	for child := title.FirstChild; child != nil; child = child.NextSibling {
		t.Logf("%#v", child)
	}
}

func TestNodeTitle(t *testing.T) {
	testx.RequireEqualf(t, "TEST", NodeTitle(root), "NodeTitle")
}
