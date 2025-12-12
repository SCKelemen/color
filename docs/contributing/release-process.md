# Release Guide

This guide explains how to release version 1.0.0 and publish to pkg.go.dev.

## Prerequisites

- Go 1.19 or later
- Git repository pushed to GitHub
- GitHub Actions enabled (for CI/CD)

## Steps to Release v1.0.0

### 1. Ensure Everything is Ready

```bash
# Run tests
go test ./...

# Build the package
go build ./...

# Check for issues
go vet ./...
```

### 2. Tag the Release

```bash
# Create and push the v1.0.0 tag
git tag v1.0.0
git push origin v1.0.0
```

### 3. Create GitHub Release (Optional)

1. Go to https://github.com/SCKelemen/color/releases
2. Click "Draft a new release"
3. Select tag `v1.0.0`
4. Title: "v1.0.0 - Initial Release"
5. Add release notes describing features
6. Click "Publish release"

### 4. Publish to pkg.go.dev

pkg.go.dev automatically indexes Go modules. To trigger indexing:

```bash
# Fetch the module (this triggers pkg.go.dev to index it)
go get github.com/SCKelemen/color@v1.0.0
```

Or visit: https://pkg.go.dev/github.com/SCKelemen/color

Note: It may take a few minutes for pkg.go.dev to index the module after the tag is pushed.

### 5. Verify

- Check GitHub Actions: https://github.com/SCKelemen/color/actions
- Check pkg.go.dev: https://pkg.go.dev/github.com/SCKelemen/color@v1.0.0
- Verify documentation renders correctly

## Future Releases

For future versions (v1.0.1, v1.1.0, etc.):

1. Update code
2. Run tests
3. Tag new version: `git tag v1.0.1 && git push origin v1.0.1`
4. pkg.go.dev will automatically index the new version

## Module Requirements for pkg.go.dev

✅ Module path matches GitHub repository  
✅ Go version specified in go.mod  
✅ Package documentation in doc.go  
✅ README.md present  
✅ LICENSE file present  
✅ Tests pass  
✅ Code compiles  

All requirements are met! 🎉

