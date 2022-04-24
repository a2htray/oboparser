package oboparser

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestParser_Execute(t *testing.T) {
	ExtendHeaderFields([]string{"ontology"})
	rd := strings.NewReader(`
format-version: 1.2
default-namespace: tripal_pub
subsetdef: MeSH_Publication_Type "MeSH Publication Types"
ontology: tpub

[Term]
id: TPUB:0000001
name: Publication Dbxref
relationship: part_of TPUB:0000037 ! Publication Details
def: A unique identifer for the publication in a remote database.  The format is a database abbreviation and a unique accession separated by a colon.  (e.g. PMID:23493063)

[Term]
id: TPUB:0000002
name: Publication
`)

	parser := New(rd, &ParserConfig{ShowNotExistField: true})
	oboObj := parser.Execute()

	formatVersion := oboObj.HeaderValue("format-version")
	assert.Equal(t, "1.2", formatVersion)
	defaultNamespace := oboObj.HeaderValue("default-namespace")
	assert.Equal(t, "tripal_pub", defaultNamespace)
	subsetDef := oboObj.HeaderValue("subsetdef")
	assert.Equal(t, "MeSH_Publication_Type \"MeSH Publication Types\"", subsetDef)
	ontology := oboObj.HeaderValue("ontology")
	assert.Equal(t, "tpub", ontology)

	stanzaNum := len(oboObj.Stanzas)
	assert.Equal(t, 2, stanzaNum)

	assert.Equal(t, "TPUB:0000001", oboObj.Stanzas[0].Values["id"])
	assert.Equal(t, "TPUB:0000002", oboObj.Stanzas[1].Values["id"])

	t.Log(oboObj.Stanzas[0].Values)
	t.Log(oboObj.Stanzas[stanzaNum-1].Values)
	t.Log(oboObj.JSON())
	t.Log(oboObj.JSONIdent("", "  "))
}

func TestParser_Execute2(t *testing.T) {
	ExtendHeaderFields([]string{"ontology"})
	rd, err := os.Open("./testData/tpub.obo")
	assert.Nil(t, err)
	parser := New(rd, &ParserConfig{ShowNotExistField: true})
	oboObj := parser.Execute()

	t.Log(oboObj.JSONIdent("", "  "))
}

func TestParser_Execute3(t *testing.T) {
	ExtendHeaderFields([]string{"ontology"})
	ExtendTypedefFields([]string{"builtin", "exact_synonym", "inverse_of_on_instance_level", "inverse_of", "instance_level_is_transitive"})

	rd, err := os.Open("./testData/relationship.obo")
	assert.Nil(t, err)
	parser := New(rd, &ParserConfig{ShowNotExistField: true})
	oboObj := parser.Execute()

	t.Log(oboObj.JSONIdent("", "  "))
}
