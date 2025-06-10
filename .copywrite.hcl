schema_version = 1

project {
  license          = ""
  copyright_holder = "Carlos De La Torre CC-BY-NC-v4 (https://creativecommons.org/licenses/by-nc/4.0/)"
  copyright_year   = "2019"
  upstream         = "elsudano/terraform-provider-vmworkstation"
  header_ignore    = [
    # internal catalog metadata (prose)
    "META.d/**/*.yaml",

    # examples used within documentation (prose)
    "examples/**",

    # GitHub issue template configuration
    ".github/ISSUE_TEMPLATE/*.yml",

    # golangci-lint tooling configuration
    ".golangci.yml",

    # GoReleaser tooling configuration
    ".goreleaser.yaml",
  ]
}
