package scvs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/interlynk-io/sbomqs/pkg/sbom"
	"github.com/samber/lo"
)

// A structured, machine readable software bill of materials (SBOM) format is present
func IsSBOMMachineReadable(d sbom.Document) bool {
	// check spec is SPDX or CycloneDX
	specs := sbom.SupportedSBOMSpecs()

	for _, spec := range specs {
		if d.Spec().GetSpecType() == spec {
			return true
		}
	}
	return false
}

// SBOM creation is automated and reproducible
func IsSBOMCreationAutomated(d sbom.Document) bool {
	if tools := d.Tools(); tools != nil {
		for _, tool := range tools {
			name := tool.GetName()
			version := tool.GetVersion()
			if name != "" && version != "" {
				return true
			}
		}
	}
	return false
}

// 2.3 Each SBOM has a unique identifier
func IsSBOMHasUniqID(d sbom.Document) bool {
	if ns := d.Spec().GetNamespace(); ns != "" {
		return true
	}
	return false
} // 2.3
func IsSBOMHasSignature(d sbom.Document) bool {
	// isSignatureExists := d.Spec().GetSignature().CheckSignatureExists()
	sig := d.Signature()

	if sig != nil {
		for _, isSignatureExists := range sig {
			return isSignatureExists.CheckSignatureExists()
		}
	}

	return false
} // 2.4

func IsSBOMSignatureCorrect(d sbom.Document) bool {
	return IsSBOMHasSignature(d)
} // 2.5
func IsSBOMSignatureVerified(d sbom.Document) bool {
	// Save signature and public key to temporary files
	signature := d.Signature()
	if signature == nil {
		return false
	}
	for _, sig := range signature {

		sigFile, err := os.CreateTemp("", "signature-*.sig")
		if err != nil {
			fmt.Println("Error creating temp file for signature:", err)
			return false
		}
		defer os.Remove(sigFile.Name())

		pubKeyFile, err := os.CreateTemp("", "publickey-*.pem")
		if err != nil {
			fmt.Println("Error creating temp file for public key:", err)
			return false
		}
		defer os.Remove(pubKeyFile.Name())

		_, err = sigFile.WriteString(sig.Value())
		if err != nil {
			fmt.Println("Error writing signature to temp file:", err)
			return false
		}
		_, err = pubKeyFile.WriteString(sig.PublicKey())
		if err != nil {
			fmt.Println("Error writing public key to temp file:", err)
			return false
		}

		// Use openssl to verify the signature
		cmd := exec.Command("openssl", "dgst", "-verify", pubKeyFile.Name(), "-signature", sigFile.Name(), "data-to-verify.txt")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error verifying signature with openssl:", err)
			return false
		}
		// // Use cosign to verify the signature
		// cmd := exec.Command("cosign", "verify-blob", "--key", pubKeyFile.Name(), "--signature", sigFile.Name(), "data-to-verify.txt")
		// output, err := cmd.CombinedOutput()
		// if err != nil {
		//     fmt.Println("Error verifying signature with cosign:", err)
		//     fmt.Println(string(output))
		//     return false
		// }

		verificationResult := strings.Contains(string(output), "Verified OK")
		fmt.Println("Verification result:", verificationResult)

		return verificationResult
	}
	return false
}

func IsSBOMTimestamped(d sbom.Document) bool {
	if d.Spec().GetCreationTimestamp() != "" {
		return true
	}
	return false
}                                                 // 2.7
func IsSBOMAnalyzedForRisk() bool                 { return false } // 2.8
func IsSBOMHasInventoryOfDependencies() bool      { return false } // 2.9
func IsSBOMInventoryContainsTestComponents() bool { return false } // 2.10
func IsSBOMHasPrimaryComponents(d sbom.Document) bool {
	//
	if d.PrimaryComponent() {
		return true
	}
	return false
}

// 2.11
func IsComponentHasIdentityID(d sbom.Document) bool {
	totalComponents := len(d.Components())
	if totalComponents == 0 {
		return false
	}

	withIdentityID := lo.CountBy(d.Components(), func(c sbom.GetComponent) bool {
		return len(c.Purls()) > 0
	})

	if totalComponents > 0 {
		if withIdentityID == totalComponents {
			return true
		} else {
			return false
		}
	}

	return false
} // 2.12

func IsComponentHasOriginID(d sbom.Document) bool {
	totalComponents := len(d.Components())
	if totalComponents == 0 {
		return false
	}

	withOriginID := lo.CountBy(d.Components(), func(c sbom.GetComponent) bool {
		return len(c.Purls()) > 0
	})

	if totalComponents > 0 {
		if withOriginID == totalComponents {
			return true
		} else {
			return false
		}
	}

	return false
} // 2.13
func IsComponentHasLicenses(d sbom.Document) bool {
	//
	totalComponents := len(d.Components())
	if totalComponents == 0 {
		return false
	}

	withLicenses := lo.CountBy(d.Components(), func(c sbom.GetComponent) bool {
		return len(c.Licenses()) > 0
	})

	if totalComponents > 0 {
		if withLicenses == totalComponents {
			return true
		} else {
			return false
		}
	}

	return false
}                                         // 2.14
func IsComponentHasVerifiedLicense() bool { return false } // 2.15
func IsComponentHasCopyright(d sbom.Document) bool {
	totalComponents := len(d.Components())
	if totalComponents == 0 {
		return false
	}

	withCopyrights := lo.CountBy(d.Components(), func(c sbom.GetComponent) bool {
		return len(c.GetCopyRight()) > 0
	})

	if totalComponents > 0 {
		if withCopyrights == totalComponents {
			return true
		} else {
			return false
		}
	}

	return false
}                                                  // 2.16
func IsComponentContainsModificationChanges() bool { return false } // 2.17
func IsComponentContainsHash() bool                { return false } // 2.18