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
	s.UniqID = "https://anchore.com/syft/dir/sbomqs-6ec18b03-96cb-4951-b299-929890c1cfc8"
	s.Organization = "interlynk"
	s.CreationTimestamp = "2023-05-04T09:33:40Z"
	s.Spdxid = "DOCUMENT"
	s.Comment = "this is a general sbom created using syft tool"
	lics := licenses.CreateCustomLicense("", "cc0-1.0")
	s.Licenses = append(s.Licenses, lics)

	var tools []sbom.GetTool
	tool := sbom.Tool{
		Name: "syft",
	}
	tools = append(tools, tool)

	pack := sbom.NewComponent()
	pack.Version = "v0.7.1"
	pack.Name = "core-js"
	pack.Spdxid = "SPDXRef-npm-core-js-3.6.5"
	pack.CopyRight = "Copyright 2001-2011 The Apache Software Foundation"
	pack.FileAnalyzed = true
	pack.ID = "Package-go-module-github.com-CycloneDX-cyclonedx-go-21b8492723f5584d"
	pack.PackageLicenseConcluded = "(LGPL-2.0-only OR LicenseRef-3)"
	pack.PackageLicenseDeclared = "(LGPL-2.0-only AND LicenseRef-3)"
	pack.DownloadLocation = "https://registry.npmjs.org/core-js/-/core-js-3.6.5.tgz"

	supplier := sbom.Supplier{
		Email: "vivekkumarsahu650@gmail.com",
	}
	pack.Supplier = supplier

	checksum := sbom.Checksum{
		Alg:     "SHA256",
		Content: "ee1300ac533cebc2d070ce3765685d5f7fca2a5a78ca15068323f68ed63d4abf",
	}

	var checksums []sbom.GetChecksum
	checksums = append(checksums, checksum)
	pack.Checksums = checksums

	extRef := sbom.ExternalReference{
		RefType: "purl",
	}
	var externalReferences []sbom.GetExternalReference
	externalReferences = append(externalReferences, extRef)
	pack.ExternalRefs = externalReferences

	var packages []sbom.GetComponent
	packages = append(packages, pack)

	doc := sbom.SpdxDoc{
		SpdxSpec:  s,
		Comps:     packages,
		SpdxTools: tools,
	}
	return doc
}

type desired struct {
	name   string
	score  float64
	result string
	key    int
	id     string
}

func TestOctSbomPass(t *testing.T) {
	doc := createDummyDocument()
	testCases := []struct {
		actual   *record
		expected desired
	}{
		{
			actual: octSpec(doc),
			expected: desired{
				name:   "octSpec",
				score:  10.0,
				result: "spdx",
				key:    SBOM_SPEC,
				id:     "SBOM Format",
			},
		},
		{
			actual: octSbomName(doc),
			expected: desired{
				name:   "octSbomName",
				score:  10.0,
				result: "nano",
				key:    SBOM_NAME,
				id:     "SPDX Elements",
			},
		},
		{
			actual: octSbomNamespace(doc),
			expected: desired{
				name:   "octSbomNamespace",
				score:  10.0,
				result: "https://anchore.com/syft/dir/sbomqs-6ec18b03-96cb-4951-b299-929890c1cfc8",
				key:    SBOM_NAMESPACE,
				id:     "SPDX Elements",
			},
		},
		{
			actual: octSbomOrganization(doc),
			expected: desired{
				name:   "octSbomOrganization",
				score:  10.0,
				result: "interlynk",
				key:    SBOM_ORG,
				id:     "SBOM Build Information",
			},
		},
		{
			actual: octSbomComment(doc),
			expected: desired{
				name:   "octSbomComment",
				score:  10.0,
				result: "this is a general sbom created using syft tool",
				key:    SBOM_COMMENT,
				id:     "SPDX Elements",
			},
		},
		{
			actual: octSbomTool(doc),
			expected: desired{
				name:   "octSbomTool",
				score:  10.0,
				result: "syft",
				key:    SBOM_TOOL,
				id:     "SBOM Build Information",
			},
		},
		{
			actual: octSbomLicense(doc),
			expected: desired{
				name:   "octSbomLicense",
				score:  10.0,
				result: "cc0-1.0",
				key:    SBOM_LICENSE,
				id:     "SPDX Elements",
			},
		},
		{
			actual: octSpecVersion(doc),
			expected: desired{
				name:   "octSpecVersion",
				score:  10.0,
				result: "SPDX-2.3",
				key:    SBOM_SPEC_VERSION,
				id:     "SPDX Elements",
			},
		},
		{
			actual: octCreatedTimestamp(doc),
			expected: desired{
				name:   "octCreatedTimestamp",
				score:  10.0,
				result: "2023-05-04T09:33:40Z",
				key:    SBOM_TIMESTAMP,
				id:     "SPDX Elements",
			},
		},
		{
			actual: octSpecSpdxID(doc),
			expected: desired{
				name:   "octSpecSpdxID",
				score:  10.0,
				result: "DOCUMENT",
				key:    SBOM_SPDXID,
				id:     "SPDX Elements",
			},
		},

		{
			actual: octMachineFormat(doc),
			expected: desired{
				name:   "octMachineFormat",
				score:  10.0,
				result: "spdx, json",
				key:    SBOM_MACHINE_FORMAT,
				id:     "Machine Readable Data Format",
			},
		},
		{
			actual: octHumanFormat(doc),
			expected: desired{
				name:   "octHumanFormat",
				score:  10.0,
				result: "json",
				key:    SBOM_HUMAN_FORMAT,
				id:     "Human Readable Data Format",
			},
		},
		{
			actual: octPackageName(doc.Components()[0]),
			expected: desired{
				name:   "octPackageName",
				score:  10.0,
				result: "core-js",
				key:    PACK_NAME,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageVersion(doc.Components()[0]),
			expected: desired{
				name:   "octPackageVersion",
				score:  10.0,
				result: "v0.7.1",
				key:    PACK_VERSION,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageSpdxID(doc.Components()[0]),
			expected: desired{
				name:   "octPackageSpdxID",
				score:  10.0,
				result: "SPDXRef-npm-core-js-3.6.5",
				key:    PACK_SPDXID,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageSupplier(doc.Components()[0]),
			expected: desired{
				name:   "octPackageSupplier",
				score:  10.0,
				result: "vivekkumarsahu650@gmail.com",
				key:    PACK_SUPPLIER,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageHash(doc.Components()[0]),
			expected: desired{
				name:   "octPackageHash",
				score:  10.0,
				result: "ee1300ac533cebc2d070ce3765685d5f7fca2a5a78ca15068323f68ed63d4abf",
				key:    PACK_HASH,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageExternalRefs(doc.Components()[0]),
			expected: desired{
				name:   "octPackageExternalRefs",
				score:  10.0,
				result: "purl:(1/1)",
				key:    PACK_EXT_REF,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageCopyright(doc.Components()[0]),
			expected: desired{
				name:   "octPackageCopyright",
				score:  10.0,
				result: "Copyright 2001-2011 The Apache Software Foundation",
				key:    PACK_COPYRIGHT,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageFileAnalyzed(doc.Components()[0]),
			expected: desired{
				name:   "octPackageFileAnalyzed",
				score:  10.0,
				result: "yes",
				key:    PACK_FILE_ANALYZED,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageConLicense(doc.Components()[0]),
			expected: desired{
				name:   "octPackageConLicense",
				score:  10.0,
				result: "(LGPL-2.0-only OR LicenseRef-3)",
				key:    PACK_LICENSE_CON,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageDecLicense(doc.Components()[0]),
			expected: desired{
				name:   "octPackageDecLicense",
				score:  10.0,
				result: "(LGPL-2.0-only AND LicenseRef-3)",
				key:    PACK_LICENSE_DEC,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageDownloadURL(doc.Components()[0]),
			expected: desired{
				name:   "octPackageDownloadURL",
				score:  10.0,
				result: "https://registry.npmjs.org/core-js/-/core-js-3.6.5.tgz",
				key:    PACK_DOWNLOAD_URL,
				id:     doc.Components()[0].GetName(),
			},
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expected.score, test.actual.score)
		assert.Equal(t, test.expected.key, test.actual.checkKey)
		assert.Equal(t, test.expected.id, test.actual.id)
		assert.Equal(t, test.expected.result, test.actual.checkValue)
	}
}

func createFailureDummyDocument() sbom.Document {
	s := sbom.NewSpec()
	s.Version = ""
	s.Format = "xml"
	s.SpecType = "cyclonedx"
	s.Name = ""
	s.UniqID = ""
	s.Organization = ""
	s.CreationTimestamp = "wrong-time-format"
	s.Spdxid = ""
	s.Comment = ""
	lics := licenses.CreateCustomLicense("", "")
	s.Licenses = append(s.Licenses, lics)

	var tools []sbom.GetTool
	tool := sbom.Tool{
		Name: "",
	}
	tools = append(tools, tool)

	pack := sbom.NewComponent()
	pack.Version = ""
	pack.Name = ""
	pack.Spdxid = ""
	pack.CopyRight = "NOASSERTION"
	pack.FileAnalyzed = false
	pack.ID = ""
	pack.PackageLicenseConcluded = "NONE"
	pack.PackageLicenseDeclared = "NOASSERTION"
	pack.DownloadLocation = ""

	supplier := sbom.Supplier{
		Email: "",
	}
	pack.Supplier = supplier

	checksum := sbom.Checksum{
		Alg:     "SHA-1",
		Content: "443238d9cf19f77ccc8cdda3ba5421ea9ea2bc78",
	}

	var checksums []sbom.GetChecksum
	checksums = append(checksums, checksum)
	pack.Checksums = checksums

	extRef := sbom.ExternalReference{
		RefType: "cpe23Type",
	}
	var externalReferences []sbom.GetExternalReference
	externalReferences = append(externalReferences, extRef)
	pack.ExternalRefs = externalReferences

	var packages []sbom.GetComponent
	packages = append(packages, pack)

	doc := sbom.SpdxDoc{
		SpdxSpec:  s,
		Comps:     packages,
		SpdxTools: tools,
	}
	return doc
}

func TestOctSbomFail(t *testing.T) {
	doc := createFailureDummyDocument()
	testCases := []struct {
		actual   *record
		expected desired
	}{
		{
			actual: octSpec(doc),
			expected: desired{
				score:  0.0,
				result: "cyclonedx",
				key:    SBOM_SPEC,
				id:     "SBOM Format",
			},
		},
		{
			actual: octSbomName(doc),
			expected: desired{
				score:  0.0,
				result: "",
				key:    SBOM_NAME,
				id:     "SPDX Elements",
			},
		},
		{
			actual: octSbomNamespace(doc),
			expected: desired{
				score:  0.0,
				result: "",
				key:    SBOM_NAMESPACE,
				id:     "SPDX Elements",
			},
		},
		{
			actual: octSbomOrganization(doc),
			expected: desired{
				score:  0.0,
				result: "",
				key:    SBOM_ORG,
				id:     "SBOM Build Information",
			},
		},
		{
			actual: octSbomComment(doc),
			expected: desired{
				score:  0.0,
				result: "",
				key:    SBOM_COMMENT,
				id:     "SPDX Elements",
			},
		},
		{
			actual: octSbomTool(doc),
			expected: desired{
				score:  0.0,
				result: "",
				key:    SBOM_TOOL,
				id:     "SBOM Build Information",
			},
		},
		{
			actual: octSbomLicense(doc),
			expected: desired{
				score:  0.0,
				result: "",
				key:    SBOM_LICENSE,
				id:     "SPDX Elements",
			},
		},
		{
			actual: octSpecVersion(doc),
			expected: desired{
				score:  0.0,
				result: "",
				key:    SBOM_SPEC_VERSION,
				id:     "SPDX Elements",
			},
		},
		{
			actual: octCreatedTimestamp(doc),
			expected: desired{
				score:  0.0,
				result: "wrong-time-format",
				key:    SBOM_TIMESTAMP,
				id:     "SPDX Elements",
			},
		},
		{
			actual: octSpecSpdxID(doc),
			expected: desired{
				score:  0.0,
				result: "",
				key:    SBOM_SPDXID,
				id:     "SPDX Elements",
			},
		},
		{
			actual: octMachineFormat(doc),
			expected: desired{
				score:  0.0,
				result: "cyclonedx, xml",
				key:    SBOM_MACHINE_FORMAT,
				id:     "Machine Readable Data Format",
			},
		},
		{
			actual: octHumanFormat(doc),
			expected: desired{
				score:  0.0,
				result: "xml",
				key:    SBOM_HUMAN_FORMAT,
				id:     "Human Readable Data Format",
			},
		},
		{
			actual: octPackageName(doc.Components()[0]),
			expected: desired{
				score:  0.0,
				result: "",
				key:    PACK_NAME,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageVersion(doc.Components()[0]),
			expected: desired{
				score:  0.0,
				result: "",
				key:    PACK_VERSION,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageSpdxID(doc.Components()[0]),
			expected: desired{
				score:  0.0,
				result: "",
				key:    PACK_SPDXID,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageSupplier(doc.Components()[0]),
			expected: desired{
				score:  0.0,
				result: "",
				key:    PACK_SUPPLIER,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageHash(doc.Components()[0]),
			expected: desired{
				score:  0.0,
				result: "",
				key:    PACK_HASH,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageExternalRefs(doc.Components()[0]),
			expected: desired{
				score:  0.0,
				result: "cpe23Type",
				key:    PACK_EXT_REF,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageCopyright(doc.Components()[0]),
			expected: desired{
				score:  0.0,
				result: "NOASSERTION",
				key:    PACK_COPYRIGHT,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageFileAnalyzed(doc.Components()[0]),
			expected: desired{
				score:  0.0,
				result: "no",
				key:    PACK_FILE_ANALYZED,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageConLicense(doc.Components()[0]),
			expected: desired{
				score:  0.0,
				result: "NONE",
				key:    PACK_LICENSE_CON,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageDecLicense(doc.Components()[0]),
			expected: desired{
				score:  0.0,
				result: "NOASSERTION",
				key:    PACK_LICENSE_DEC,
				id:     doc.Components()[0].GetName(),
			},
		},
		{
			actual: octPackageDownloadURL(doc.Components()[0]),
			expected: desired{
				score:  0.0,
				result: "",
				key:    PACK_DOWNLOAD_URL,
				id:     doc.Components()[0].GetName(),
			},
		},
	}

	for _, test := range testCases {
		assert.Equal(t, test.expected.score, test.actual.score)
		assert.Equal(t, test.expected.key, test.actual.checkKey)
		assert.Equal(t, test.expected.id, test.actual.id)
		assert.Equal(t, test.expected.result, test.actual.checkValue)
	}
}
