package bundle

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveToLocal_GitHubRelease(t *testing.T) {
	t.Run("resolves tagged release and verifies checksum", func(t *testing.T) {
		archiveBytes := buildBundleArchive(t, archiveFixture{tag: "v1.2.3", bundleVersion: "v1.2.3", presetName: "mixed"})
		checksum := sha256Hex(archiveBytes)

		server := newGitHubReleaseTestServer(t, githubReleaseServerFixture{
			repo:         "qbicsoftware/opencode-config-bundle",
			tag:          "v1.2.3",
			archiveBytes: archiveBytes,
			checksums:    fmt.Sprintf("%s  opencode-config-bundle-v1.2.3.tar.gz\n", checksum),
		})
		defer server.Close()

		oldClient := githubHTTPClient
		oldAPIBaseURL := githubAPIBaseURL
		oldDownloadBaseURL := githubDownloadBaseURL
		oldCacheDir := githubCacheDirOverride
		githubHTTPClient = server.Client()
		githubAPIBaseURL = server.URL
		githubDownloadBaseURL = server.URL
		githubCacheDirOverride = t.TempDir()
		defer func() {
			githubHTTPClient = oldClient
			githubAPIBaseURL = oldAPIBaseURL
			githubDownloadBaseURL = oldDownloadBaseURL
			githubCacheDirOverride = oldCacheDir
		}()

		bundleRoot, cleanup, err := ResolveToLocal("github-release", "qbicsoftware/opencode-config-bundle", "v1.2.3")
		if err != nil {
			t.Fatalf("ResolveToLocal() error = %v", err)
		}
		defer cleanup()

		manifestPath := filepath.Join(bundleRoot, "opencode-bundle.manifest.json")
		if _, err := os.Stat(manifestPath); err != nil {
			t.Fatalf("expected manifest at %s: %v", manifestPath, err)
		}
	})

	t.Run("fails when checksums do not match", func(t *testing.T) {
		archiveBytes := buildBundleArchive(t, archiveFixture{tag: "v1.2.3", bundleVersion: "v1.2.3", presetName: "mixed"})

		server := newGitHubReleaseTestServer(t, githubReleaseServerFixture{
			repo:         "qbicsoftware/opencode-config-bundle",
			tag:          "v1.2.3",
			archiveBytes: archiveBytes,
			checksums:    "deadbeef  opencode-config-bundle-v1.2.3.tar.gz\n",
		})
		defer server.Close()

		oldClient := githubHTTPClient
		oldAPIBaseURL := githubAPIBaseURL
		oldDownloadBaseURL := githubDownloadBaseURL
		oldCacheDir := githubCacheDirOverride
		githubHTTPClient = server.Client()
		githubAPIBaseURL = server.URL
		githubDownloadBaseURL = server.URL
		githubCacheDirOverride = t.TempDir()
		defer func() {
			githubHTTPClient = oldClient
			githubAPIBaseURL = oldAPIBaseURL
			githubDownloadBaseURL = oldDownloadBaseURL
			githubCacheDirOverride = oldCacheDir
		}()

		_, cleanup, err := ResolveToLocal("github-release", "qbicsoftware/opencode-config-bundle", "v1.2.3")
		if cleanup != nil {
			defer cleanup()
		}
		if err == nil {
			t.Fatalf("ResolveToLocal() error = nil, want checksum failure")
		}
		if !strings.Contains(err.Error(), "SHA256 mismatch") {
			t.Fatalf("ResolveToLocal() error = %v, want SHA256 mismatch", err)
		}
	})

	t.Run("fails when release has no tarball asset", func(t *testing.T) {
		server := newGitHubReleaseTestServer(t, githubReleaseServerFixture{
			repo:         "qbicsoftware/opencode-config-bundle",
			tag:          "v1.2.3",
			archiveBytes: nil,
		})
		defer server.Close()

		oldClient := githubHTTPClient
		oldAPIBaseURL := githubAPIBaseURL
		oldDownloadBaseURL := githubDownloadBaseURL
		oldCacheDir := githubCacheDirOverride
		githubHTTPClient = server.Client()
		githubAPIBaseURL = server.URL
		githubDownloadBaseURL = server.URL
		githubCacheDirOverride = t.TempDir()
		defer func() {
			githubHTTPClient = oldClient
			githubAPIBaseURL = oldAPIBaseURL
			githubDownloadBaseURL = oldDownloadBaseURL
			githubCacheDirOverride = oldCacheDir
		}()

		_, cleanup, err := ResolveToLocal("github-release", "qbicsoftware/opencode-config-bundle", "v1.2.3")
		if cleanup != nil {
			defer cleanup()
		}
		if err == nil {
			t.Fatalf("ResolveToLocal() error = nil, want missing asset failure")
		}
		if !strings.Contains(err.Error(), "no bundle asset") {
			t.Fatalf("ResolveToLocal() error = %v, want no bundle asset message", err)
		}
	})
}

type archiveFixture struct {
	tag           string
	bundleVersion string
	presetName    string
}

func buildBundleArchive(t *testing.T, fixture archiveFixture) []byte {
	t.Helper()

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)

	root := fmt.Sprintf("opencode-config-bundle-%s", fixture.tag)
	manifest := fmt.Sprintf(`{
		"manifest_version": "1.0.0",
		"bundle_name": "test-bundle",
		"bundle_version": %q,
		"presets": [{"name": %q, "entrypoint": "opencode.json"}]
	}`, fixture.bundleVersion, fixture.presetName)

	writeTarFile(t, tw, root+"/opencode-bundle.manifest.json", []byte(manifest))
	writeTarFile(t, tw, root+"/opencode.json", []byte(`{"agents":[]}`))

	if err := tw.Close(); err != nil {
		t.Fatalf("closing tar writer: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("closing gzip writer: %v", err)
	}

	return buf.Bytes()
}

func writeTarFile(t *testing.T, tw *tar.Writer, name string, body []byte) {
	t.Helper()
	if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0644, Size: int64(len(body))}); err != nil {
		t.Fatalf("WriteHeader(%s): %v", name, err)
	}
	if _, err := tw.Write(body); err != nil {
		t.Fatalf("Write(%s): %v", name, err)
	}
}

type githubReleaseServerFixture struct {
	repo         string
	tag          string
	archiveBytes []byte
	checksums    string
}

func newGitHubReleaseTestServer(t *testing.T, fixture githubReleaseServerFixture) *httptest.Server {
	t.Helper()

	archiveName := fmt.Sprintf("opencode-config-bundle-%s.tar.gz", fixture.tag)
	handler := http.NewServeMux()

	handler.HandleFunc("/repos/"+fixture.repo+"/releases/tags/"+fixture.tag, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(buildReleaseJSON("http://"+r.Host, fixture, archiveName)))
	})
	handler.HandleFunc("/repos/"+fixture.repo+"/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(buildReleaseJSON("http://"+r.Host, fixture, archiveName)))
	})
	handler.HandleFunc("/downloads/"+fixture.repo+"/releases/download/"+fixture.tag+"/"+archiveName, func(w http.ResponseWriter, r *http.Request) {
		if len(fixture.archiveBytes) == 0 {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/gzip")
		_, _ = w.Write(fixture.archiveBytes)
	})
	handler.HandleFunc("/downloads/"+fixture.repo+"/releases/download/"+fixture.tag+"/opencode-config-bundle-"+fixture.tag+"-checksums.txt", func(w http.ResponseWriter, r *http.Request) {
		if fixture.checksums == "" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(fixture.checksums))
	})

	return httptest.NewServer(handler)
}

func buildReleaseJSON(baseURL string, fixture githubReleaseServerFixture, archiveName string) string {
	if len(fixture.archiveBytes) == 0 {
		return fmt.Sprintf(`{"tag_name":%q,"assets":[]}`, fixture.tag)
	}

	return fmt.Sprintf(`{
		"tag_name": %q,
		"assets": [
			{"name": %q, "browser_download_url": %q},
			{"name": %q, "browser_download_url": %q}
		]
	}`,
		fixture.tag,
		archiveName,
		fmt.Sprintf("%s/downloads/%s/releases/download/%s/%s", baseURL, fixture.repo, fixture.tag, archiveName),
		"opencode-config-bundle-"+fixture.tag+"-checksums.txt",
		fmt.Sprintf("%s/downloads/%s/releases/download/%s/%s", baseURL, fixture.repo, fixture.tag, "opencode-config-bundle-"+fixture.tag+"-checksums.txt"),
	)
}
