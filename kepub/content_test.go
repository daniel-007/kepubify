package kepub

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestCleanHTML(t *testing.T) {
	h := `<meta  content="urn:uuid:asd--asdasd-asdasdas-dasdasd234234"   name="Adept.expected.resource"   />��<st1:asd></st1:asd><o:p></o:p><h1></h1><h3></h3><h2>test</h2><style></style>`
	cleanHTML(&h)
	assert.Equal(t, " <h2>test</h2><style type=\"text/css\"></style>", h, "should be equal if cleaned correctly")
}

func TestSmartenPunctuation(t *testing.T) {
	h := `-- --- <!--test-->`
	smartenPunctuation(&h)
	assert.Equal(t, " &#x2014;   &#x2013;  <!-- test -->", h, "should be equal if smartened correctly")
}

func TestOpenSelfClosingPs(t *testing.T) {
	h := `<p>test</p><p /><p  /><p>test</p>`
	openSelfClosingPs(&h)
	assert.Equal(t, "<p>test</p><p></p><p></p><p>test</p>", h, "should be equal if reopened correctly")
}

func TestAddSpans(t *testing.T) {
	h := `<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <title>Test Book 1</title>
    <meta content="http://www.w3.org/1999/xhtml; charset=utf-8" http-equiv="Content-Type"/>
    <style type="text/css">
        @page { margin-bottom: 5.000000pt; margin-top: 5.000000pt; }
    </style>
</head>
<body id="p1">
	<p>This is the first sentence. This is the second sentence. This is the third sentence.</p>
	<p>This is the first sentence. This is the second sentence? This is the third sentence!</p>
	<p>This is the first <b>sentence</b>. This is the second sentence? This is the third sentence!</p>
	<p>This is the first <b>sentence. This is the </b>second sentence? This is the third sentence!</p>
	<p>This is <i>t<b>h</b>e</i> first <a href="test.html">sentence <b>here</b></a>. This is the second sentence? This is the third sentence!</p>
	<ul>
		<li>test</li>
		<li>test</li>
	</ul>
    </div>
</body>
</html>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(h))
	assert.Nil(t, err, "err should be nil")

	addSpans(doc)

	nh, err := doc.Html()
	assert.Nil(t, err, "err should be nil")

	nh = html.UnescapeString(nh)

	hs := sha256.New()
	hs.Write([]byte(nh))

	hxs := fmt.Sprintf("%x", hs.Sum(nil))

	assert.Equal(t, "ae78fe3c38c263e2ad43879cb6c2eaf3c0e8dafdf1fdf877bff449f2c2c44eee", hxs, "hash of content should be equal if processed correctly")
}

func TestAddDivs(t *testing.T) {
	h := `<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <title>Test Book 1</title>
    <meta content="http://www.w3.org/1999/xhtml; charset=utf-8" http-equiv="Content-Type"/>
    <style type="text/css">
        @page { margin-bottom: 5.000000pt; margin-top: 5.000000pt; }
    </style>
</head>
<body id="p1">
	<p>This is the first sentence. This is the second sentence. This is the third sentence.</p>
	<p>This is the first sentence. This is the second sentence? This is the third sentence!</p>
	<p>This is the first <b>sentence</b>. This is the second sentence? This is the third sentence!</p>
	<p>This is the first <b>sentence. This is the </b>second sentence? This is the third sentence!</p>
	<p>This is <i>t<b>h</b>e</i> first <a href="test.html">sentence <b>here</b></a>. This is the second sentence? This is the third sentence!</p>
	<ul>
		<li>test</li>
		<li>test</li>
	</ul>
    </div>
</body>
</html>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(h))
	assert.Nil(t, err, "err should be nil")

	addDivs(doc)

	nh, err := doc.Html()
	assert.Nil(t, err, "err should be nil")

	hs := sha256.New()
	hs.Write([]byte(nh))

	hxs := fmt.Sprintf("%x", hs.Sum(nil))

	assert.Equal(t, "d24ae5a8f438358828d50b036007fe06c9e24b55d6aa238f4628a24d77a15485", hxs, "hash of content should be equal if processed correctly")
}

func TestProcess(t *testing.T) {
	h := `<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <title>Test Book 1</title>
    <meta content="http://www.w3.org/1999/xhtml; charset=utf-8" http-equiv="Content-Type"/>
    <style type="text/css">
        @page { margin-bottom: 5.000000pt; margin-top: 5.000000pt; }
    </style>
</head>
<body id="p1">
	<p>This is the first sentence. This is the second sentence. This is the third sentence.</p>
	<p>This is the first sentence. This is the second sentence? This is the third sentence!</p>
	<p>This is the first <b>sentence</b>. This is the second sentence? This is the third sentence!</p>
	<p>This is the first <b>sentence. This is the </b>second sentence? This is the third sentence!</p>
	<p>This is <i>t<b>h</b>e</i> first <a href="test.html">sentence <b>here</b></a>. This is the second sentence? This is the third sentence!</p>
	<ul>
		<li>test</li>
		<li>test</li>
	</ul>
    </div>
</body>
</html>`

	process(&h)

	hs := sha256.New()
	hs.Write([]byte(h))

	hxs := fmt.Sprintf("%x", hs.Sum(nil))

	assert.Equal(t, "3abc0810906b322e3860b3d7fd1bafd5133a4a66ced286497eaafb40c94612fd", hxs, "hash of content should be equal if processed correctly")
}
