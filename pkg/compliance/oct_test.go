package compliance

import (
	"testing"

	"github.com/interlynk-io/sbomqs/pkg/licenses"
	"github.com/interlynk-io/sbomqs/pkg/sbom"
	"gotest.tools/assert"
)

func createDummyDocument() sbom.Document {
	s := sbom.NewSpec()
	s.Version = "SPDX-2.3"
	s.Format = "json"
	s.SpecType = "spdx"
	s.Name = "nano"
	s.Namespace = "https://anchore.com/syft/dir/sbomqs-6ec18b03-96cb-4951-b299-929890c1cfc8"
	s.Organization = "interlynk"
	s.CreationTimestamp = "2023-05-04T09:33:40Z"
	s.Spdxid = "DOCUMENT"
	s.Comment = "this is a general sbom created using syft tool"
	lics := licenses.CreateCustomLicense("", "cc0-1.0")
	s.Licenses = append(s.Licenses, lics)

	pack := sbom.NewComponent()
	pack.Version = "v0.7.1"
	pack.Name = "core-js"
	pack.Spdxid = "SPDXRef-npm-core-js-3.6.5"
	pack.CopyRight = "Copyright 2001-2011 The Apache Software Foundation"
	pack.FileAnalyzed = true
	pack.Id = "Package-go-module-github.com-CycloneDX-cyclonedx-go-21b8492723f5584d"
	pack.PackageLicenseConcluded = "(LGPL-2.0-only OR LicenseRef-3)"
	pack.PackageLicenseDeclared = "(LGPL-2.0-only AND LicenseRef-3)"
	pack.DownloadLocation = "https://registry.npmjs.org/core-js/-/core-js-3.6.5.tgz"

	var packages []sbom.GetComponent
	packages = append(packages, pack)

	doc := sbom.SpdxDoc{
		SpdxSpec: s,
		Comps:    packages,
	}
	return doc
}

func TestOctSbomSuccessful(t *testing.T) {
	doc := createDummyDocument()
	type desired struct {
		score  float64
		result string
		key    int
		id     string
	}
	testCases := []struct {
		actual   *record
		expected desired
	}{
		{
			actual: octSpec(doc),
			expected: desired{
				score:  10.0,
				result: "spdx",
				key:    SBOM_SPEC,
				id:     "SBOM DataFormat",
			},
		},
		{
			actual: octSbomName(doc),
			expected: desired{
				score:  10.0,
				result: "nano",
				key:    SBOM_NAME,
				id:     "doc",
			},
		},
		{
			actual: octSbomNamespace(doc),
			expected: desired{
				score:  10.0,
				result: "https://anchore.com/syft/dir/sbomqs-6ec18b03-96cb-4951-b299-929890c1cfc8",
				key:    SBOM_NAMESPACE,
				id:     "doc",
			},
		},
		{
			actual: octSbomOrganization(doc),
			expected: desired{
				score:  10.0,
				result: "interlynk",
				key:    SBOM_ORG,
				id:     "SBOM Build Information",
			},
		},
		{
			actual: octSbomComment(doc),
			expected: desired{
				score:  10.0,
				result: "this is a general sbom created using syft tool",
				key:    SBOM_COMMENT,
				id:     "doc",
			},
		},
		{
			actual: octSbomLicense(doc),
			expected: desired{
				score:  10.0,
				result: "cc0-1.0",
				key:    SBOM_LICENSE,
				id:     "doc",
			},
		},
		{
			actual: octSpecVersion(doc),
			expected: desired{
				score:  10.0,
				result: "SPDX-2.3",
				key:    SBOM_SPEC_VERSION,
				id:     "doc",
			},
		},
		{
			actual: octCreatedTimestamp(doc),
			expected: desired{
				score:  10.0,
				result: "2023-05-04T09:33:40Z",
				key:    SBOM_TIMESTAMP,
				id:     "doc",
			},
		},
		{
			actual: octSpecSpdxID(doc),
			expected: desired{
				score:  10.0,
				result: "DOCUMENT",
				key:    SBOM_SPDXID,
				id:     "doc",
			},
		},

		{
			actual: octMachineFormat(doc),
			expected: desired{
				score:  10.0,
				result: "spdx, json",
				key:    SBOM_MACHINE_FORMAT,
				id:     "Machine Readable Data Format",
			},
		},
		{
			actual: octHumanFormat(doc),
			expected: desired{
				score:  10.0,
				result: "json",
				key:    SBOM_HUMAN_FORMAT,
				id:     "Human Readable Data Format",
			},
		},
		{
			actual: octPackageName(doc.Components()[0]),
			expected: desired{
				score:  10.0,
				result: "core-js",
				key:    PACK_NAME,
				id:     doc.Components()[0].GetID(),
			},
		},
		{
			actual: octPackageVersion(doc.Components()[0]),
			expected: desired{
				score:  10.0,
				result: "v0.7.1",
				key:    PACK_VERSION,
				id:     doc.Components()[0].GetID(),
			},
		},
		{
			actual: octPackageSpdxID(doc.Components()[0]),
			expected: desired{
				score:  10.0,
				result: "SPDXRef-npm-core-js-3.6.5",
				key:    PACK_SPDXID,
				id:     doc.Components()[0].GetID(),
			},
		},
		{
			actual: octPackageCopyright(doc.Components()[0]),
			expected: desired{
				score:  10.0,
				result: "Copyright 2001-2011 The Apache Software Foundation",
				key:    PACK_COPYRIGHT,
				id:     doc.Components()[0].GetID(),
			},
		},
		{
			actual: octPackageFileAnalyzed(doc.Components()[0]),
			expected: desired{
				score:  10.0,
				result: "yes",
				key:    PACK_FILE_ANALYZED,
				id:     doc.Components()[0].GetID(),
			},
		},
		{
			actual: octPackageConLicense(doc.Components()[0]),
			expected: desired{
				score:  10.0,
				result: "(LGPL-2.0-only OR LicenseRef-3)",
				key:    PACK_LICENSE_CON,
				id:     doc.Components()[0].GetID(),
			},
		},
		{
			actual: octPackageDecLicense(doc.Components()[0]),
			expected: desired{
				score:  10.0,
				result: "(LGPL-2.0-only AND LicenseRef-3)",
				key:    PACK_LICENSE_DEC,
				id:     doc.Components()[0].GetID(),
			},
		},
		{
			actual: octPackageDownloadUrl(doc.Components()[0]),
			expected: desired{
				score:  10.0,
				result: "https://registry.npmjs.org/core-js/-/core-js-3.6.5.tgz",
				key:    PACK_DOWNLOAD_URL,
				id:     doc.Components()[0].GetID(),
			},
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expected.score, test.actual.score)
		assert.Equal(t, test.expected.key, test.actual.check_key)
		assert.Equal(t, test.expected.id, test.actual.id)
		assert.Equal(t, test.expected.result, test.actual.check_value)
	}
}
