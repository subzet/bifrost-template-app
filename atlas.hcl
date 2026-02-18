data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./model",   # relative path to the package with your GORM structs
    "--dialect", "sqlite"            # use "sqlite" for local + Turso/libSQL compatibility
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url   # ← this references the loaded schema (no "gorm://" here!)

  dev = "docker://sqlite/3/dev?search_path=main"   # for accurate diff simulation

  migration {
    dir = "file://migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}