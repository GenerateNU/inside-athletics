data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/models",
    "--dialect", "postgres", // | postgres | sqlite | sqlserver
  ]
}

env "dev" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/15/dev"
  url = getenv("DEV_DB_CONNECTION_STRING")
  schemas = ["public"]
  migration {
    dir = "file://internal/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}


env "prod" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/15/dev"
  url = getenv("PROD_MIGRATION_DB_CONNECTION_STRING")
  schemas = ["public"]
  migration {
    dir = "file://internal/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}