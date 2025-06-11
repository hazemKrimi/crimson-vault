data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/types",
    "--dialect", "sqlite",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "sqlite://file?mode=memory&_fk=1"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
